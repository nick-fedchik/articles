"""On-demand execution records, not a scheduler or scientific certification."""

import argparse
import hashlib
import json
import re
from datetime import datetime, timedelta, timezone
from pathlib import Path, PurePosixPath

ROOT = Path(__file__).resolve().parents[1]
BASELINE = "804d4772a0ea239a5ccd9e08f19f6a321ceb9a3a"
TASK_IDS = ("OPS-01", "MAP-01", "MAP-02", "MAP-03", "MAP-04")
STAGES = ("Inventory and execution control", "Full-book map", "Scientific review",
          "Cross-chapter consistency", "Substantive corrections", "Editorial revision",
          "Formatting and validation", "Independent review and release decision")
STATES = {"pending", "running", "paused", "blocked", "completed"}


def require(condition, message):
    if not condition:
        raise ValueError(message)


def utc(value):
    require(isinstance(value, str) and re.fullmatch(r"\d{4}-\d\d-\d\dT\d\d:\d\d:\d\dZ", value),
            f"malformed UTC timestamp: {value!r}")
    return datetime.strptime(value, "%Y-%m-%dT%H:%M:%SZ").replace(tzinfo=timezone.utc)


def file_path(root, value, output=False):
    require(isinstance(value, str), "path must be a string")
    path = PurePosixPath(value)
    require(value and not path.is_absolute() and str(path) == value
            and ".." not in path.parts and not re.search(r"[\s\\:]", value),
            f"malformed path: {value!r}")
    resolved = (root / path).resolve(strict=True)
    require(resolved.is_relative_to(root.resolve()) and resolved.is_file(), f"unsafe path: {value}")
    if output:
        require(path.parts[0] == "review" and path.parts[1:2] != ("tmp",)
                and resolved.is_relative_to(root / "review")
                and not resolved.is_relative_to(root / "review/tmp"),
                f"permanent output must be in review outside review/tmp: {value}")
    return resolved


def sha256(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def verify_hash(root, record, kind):
    require(isinstance(record["sha256"], str) and re.fullmatch(r"[0-9a-f]{64}", record["sha256"]),
            f"malformed SHA256: {record['path']}")
    require(sha256(file_path(root, record["path"], kind == "output")) == record["sha256"],
            f"{kind} hash mismatch (stale or changed artifact): {record['path']}")


def inventory(root):
    chapters, companions = [], []
    for path in sorted((root / "ExpertSystem").glob("*.md")):
        if path.name == "AGENTS.md":
            continue
        relative = path.relative_to(root).as_posix()
        content = file_path(root, relative).read_bytes()
        title = next((line[2:].strip() for line in content.decode("utf-8-sig").splitlines()
                      if line.startswith("# ")), "")
        require(title, f"missing H1 title: {relative}")
        record = {"path": relative, "title": title, "sha256": hashlib.sha256(content).hexdigest()}
        match = re.fullmatch(r"ch(\d{2})-.+\.md", path.name)
        if match:
            record.update(number=int(match[1]), scientific_review="pending")
            chapters.append(record)
        else:
            companions.append(record)
    require([record["number"] for record in chapters] == list(range(1, 41)),
            "source chapter coverage must be exactly ch01 through ch40, once each")
    return chapters, companions


def scope(index):
    return list(range(1, 41)) if index == 0 else list(range(10 * index - 9, 10 * index + 1))


def new_ledger(root, now=None):
    now = (now or datetime.now(timezone.utc)).replace(microsecond=0)
    stamp = now.strftime("%Y-%m-%dT%H:%M:%SZ")
    chapters, companions = inventory(root)
    tasks = []
    for index, task_id in enumerate(TASK_IDS):
        inputs = [record for record in chapters if record["number"] in scope(index)]
        if index == 0:
            inputs += companions
        tasks.append({"id": task_id, "stage": "S1" if index == 0 else "S2",
                      "chapters": scope(index), "owner": "control-agent" if index == 0 else f"map-agent-{index}",
                      "status": "running", "checkpoint_at": stamp,
                      "deadline_at": (now + timedelta(minutes=10)).strftime("%Y-%m-%dT%H:%M:%SZ"),
                      "planned": "Control setup" if index == 0 else f"Map chapters {scope(index)[0]}-{scope(index)[-1]}",
                      "doing": "Registered; execution and main-agent acceptance remain unverified.",
                      "findings": "No scientific findings recorded.",
                      "next_action": "Main reviewer: verify scope and permanent artifacts before recording acceptance.",
                      "inputs": [{"path": record["path"], "sha256": record["sha256"]} for record in inputs],
                      "outputs": [], "acceptance": None})
    return {"schema_version": 1, "baseline_commit": BASELINE, "created_at": stamp, "updated_at": stamp,
            "checkpoint_minutes": 10, "chapters": chapters, "companions": companions, "tasks": tasks,
            "stages": [{"id": f"S{index}", "status": "running" if index == 1 else "pending",
                        "planned": name, "next_action": "Meet PROGRAM.md acceptance criteria; main reviewer verifies."}
                       for index, name in enumerate(STAGES, 1)]}


def load_ledger(path):
    def unique(pairs):
        result = {}
        for key, value in pairs:
            require(key not in result, f"duplicate JSON key: {key}")
            result[key] = value
        return result
    return json.loads(path.read_text(encoding="utf-8"), object_pairs_hook=unique)


def check(root, data, now=None):
    now = now or datetime.now(timezone.utc)
    issues = []
    try:
        require(data["schema_version"] == 1 and data["baseline_commit"] == BASELINE, "wrong schema or baseline")
        created, updated = utc(data["created_at"]), utc(data["updated_at"])
        require(created <= updated <= now, "invalid ledger timestamp order or future timestamp")
        require(data["checkpoint_minutes"] == 10, "checkpoint threshold must be 10 minutes")
        require([stage["id"] for stage in data["stages"]] == [f"S{index}" for index in range(1, 9)],
                "stages must be exactly S1-S8")
        require(all(stage["status"] in STATES for stage in data["stages"]), "invalid stage status")
        chapters, companions = inventory(root)
        require([record["number"] for record in data["chapters"]] == list(range(1, 41)),
                "ledger chapter coverage must be exactly ch01 through ch40, once each")
        for key, live in (("chapters", chapters), ("companions", companions)):
            require([record["path"] for record in data[key]] == [record["path"] for record in live],
                    f"{key} paths or inventory changed")
            for record, current in zip(data[key], live):
                verify_hash(root, record, "source snapshot")
                require(record["title"] == current["title"], f"title snapshot changed: {record['path']}")
        require(all(record["scientific_review"] in STATES | {"stale"} for record in data["chapters"]),
                "invalid scientific review status")
        tasks = {task["id"]: task for task in data["tasks"]}
        require(len(tasks) == len(data["tasks"]), "duplicate task IDs")
        for index, task_id in enumerate(TASK_IDS):
            require(task_id in tasks and tasks[task_id]["chapters"] == scope(index)
                    and tasks[task_id]["stage"] == ("S1" if index == 0 else "S2"),
                    f"missing task or incorrect coverage/stage: {task_id}")
    except (ValueError, TypeError, KeyError, OSError, AttributeError) as error:
        return [f"ERROR: {error}"]
    for task in data["tasks"]:
        try:
            require(re.fullmatch(r"[A-Z]+-\d{2}", task["id"]) and task["status"] in STATES
                    and task["stage"] in {f"S{index}" for index in range(1, 9)}, "invalid task ID/status/stage")
            require(all(isinstance(task[key], str) and task[key].strip()
                        for key in ("owner", "planned", "doing", "findings", "next_action")), "empty task details")
            numbers = task["chapters"]
            require(numbers and all(type(number) is int and 1 <= number <= 40 for number in numbers)
                    and len(numbers) == len(set(numbers)), "invalid task chapter coverage")
            checkpoint, deadline = utc(task["checkpoint_at"]), utc(task["deadline_at"])
            require(created <= checkpoint <= updated and deadline == checkpoint + timedelta(minutes=10),
                    "invalid checkpoint/deadline order or threshold")
            for kind in ("inputs", "outputs"):
                paths = [record["path"] for record in task[kind]]
                require(len(paths) == len(set(paths)), f"duplicate {kind} paths")
                for record in task[kind]:
                    verify_hash(root, record, "input" if kind == "inputs" else "output")
            required = {record["path"] for record in chapters if record["number"] in numbers}
            require(required <= {record["path"] for record in task["inputs"]}, "missing chapter input hashes")
            if task["status"] == "completed":
                require(task["outputs"], "completed task has no permanent output")
                acceptance = task["acceptance"]
                require(isinstance(acceptance, dict) and all(isinstance(acceptance.get(key), str)
                        and acceptance[key].strip() for key in ("by", "at", "evidence")), "missing acceptance evidence")
                require(acceptance["by"] != task["owner"] and created <= utc(acceptance["at"]) <= updated,
                        "invalid acceptance reviewer or timestamp")
            if task["status"] == "running" and now > deadline:
                issues.append(f"OVERDUE: {task['id']}; suspicion only, not a process diagnosis; check owner/state.")
        except (ValueError, TypeError, KeyError, OSError, AttributeError) as error:
            issues.append(f"ERROR: {task['id']}: {error}")
    for stage in data["stages"]:
        if stage["status"] == "completed":
            stage_tasks = [task for task in data["tasks"] if task["stage"] == stage["id"]]
            coverage = {number for task in stage_tasks for number in task["chapters"]}
            if not stage_tasks or any(task["status"] != "completed" for task in stage_tasks) or coverage != set(range(1, 41)):
                issues.append(f"ERROR: {stage['id']}: stage completion requires accepted tasks covering all 40 chapters")
            if stage["id"] == "S3" and any(record["scientific_review"] != "completed" for record in data["chapters"]):
                issues.append("ERROR: S3: scientific review is incomplete")
    for record in data["chapters"]:
        if record["scientific_review"] == "completed" and not any(
                task["stage"] == "S3" and task["status"] == "completed" and record["number"] in task["chapters"]
                for task in data["tasks"]):
            issues.append(f"ERROR: chapter {record['number']:02}: scientific completion requires an accepted S3 task")
    return issues


def report(root, data, issues):
    def cell(value):
        return str(value).replace("|", "&#124;").replace("\n", " ")
    lines = ["# Execution Progress", "", f"Generated UTC: {datetime.now(timezone.utc).strftime('%Y-%m-%dT%H:%M:%SZ')}",
             f"Baseline: `{data['baseline_commit']}`. One process: chapters 01-40.", "",
             "On-demand check only; no background watchdog. Overdue is suspicion, not a process diagnosis.",
             "Counts describe recorded execution, not scientific truth or verified book quality.",
             f"Registered: {len(data['chapters'])}/40 chapter records; {len(data['companions'])} companions; "
             f"{len(data['tasks'])} tasks; {sum(task['status'] == 'completed' for task in data['tasks'])} recorded completed.",
             "", "## Findings", "", *(issues or ["No execution-integrity issues detected. Scientific review remains unverified."]),
             "", "## Stages", "", "| Stage | State | Planned | Next Action |", "|---|---|---|---|"]
    lines += ["| " + " | ".join(cell(stage[key]) for key in ("id", "status", "planned", "next_action")) + " |"
              for stage in data["stages"]]
    lines += ["", "## Tasks", "", "| Task | Owner | State | Planned | Doing | Findings | Next Action |",
              "|---|---|---|---|---|---|---|"]
    lines += ["| " + " | ".join(cell(task[key]) for key in
              ("id", "owner", "status", "planned", "doing", "findings", "next_action")) + " |" for task in data["tasks"]]
    lines += ["", "## All 40 Chapters", "", "| Chapter / Title | SHA256 Snapshot | Scientific Review (Recorded) |", "|---|---|---|"]
    chapters = {record["number"]: record for record in data["chapters"]}
    for number in range(1, 41):
        record = chapters.get(number)
        lines.append(f"| {number:02}: [{cell(record['title'])}](../{record['path']}) | `{record['sha256']}` | "
                     f"{cell(record['scientific_review'])} |" if record else f"| {number:02}: MISSING | - | pending |")
    lines += ["", "## Companion Inventory", ""]
    lines += [f"- [{cell(record['title'])}](../{record['path']})" for record in data["companions"]]
    active = [task for task in data["tasks"] if task["status"] not in {"completed", "paused"}]
    resume = next((task for task in active if task["id"] == data.get("resume_next_task")), None)
    next_action = (resume or active[0])["next_action"] if active else "Register the next bounded task from PROGRAM.md; no scientific completion is inferred."
    lines += ["", "## Next Action", "", next_action,
              "Record permanent output path/SHA256 and acceptance.by/at/evidence in execution.json only after review.",
              "Keep task inputs separate from outputs. Changed inputs invalidate old results even after inventory refresh.",
              "Run `PYTHONDONTWRITEBYTECODE=1 python3 review/check_execution.py report` from the repository root.", ""]
    (root / "review/PROGRESS.md").write_text("\n".join(lines), encoding="utf-8")


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("command", choices=("init", "check", "report"))
    command = parser.parse_args().command
    try:
        path = ROOT / "review/execution.json"
        if command == "init":
            data = new_ledger(ROOT)
            with path.open("x", encoding="utf-8") as output:
                json.dump(data, output, ensure_ascii=False, indent=2)
                output.write("\n")
        else:
            data = load_ledger(path)
        issues = check(ROOT, data)
        if command == "report":
            report(ROOT, data, issues)
        print("\n".join(issues) or "PASS: execution integrity; all 40 chapters inventoried, not scientifically certified.")
        return int(bool(issues))
    except (ValueError, TypeError, KeyError, OSError, AttributeError) as error:
        print(f"ERROR: {error}")
        return 1


if __name__ == "__main__":
    raise SystemExit(main())
