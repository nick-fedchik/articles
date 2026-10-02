"""Build Ukrainian reading editions with Pandoc and validate their structure."""

import json
import posixpath
import re
import shutil
import subprocess
from html.parser import HTMLParser
from pathlib import Path
from urllib.parse import unquote, urlsplit
import xml.etree.ElementTree as ET
import zipfile


ROOT = Path(__file__).resolve().parent
TITLE = "Від модема до автономного руху"
AUTHOR = "Микола Федчик / Mykola Fedchyk"
STEM = "vid-modema-do-avtonomnoho-rukhu"
SOURCE_CODE = re.compile(r"(?<![A-Z0-9])[HNAPVDMSRU]\d{2}(?![A-Z0-9])")
EDITORIAL_MARKERS = re.compile(
    r"—|екосистем|ecosystem|чесн|Таким чином|контур|фундаментальн|"
    r"унікальн|бездоганн|вичерпн|потужн|революційн|AutoForge|AFES",
    re.IGNORECASE,
)


def require(condition, message):
    if not condition:
        raise ValueError(message)


def pandoc(*arguments, input_text=None):
    result = subprocess.run(
        ["pandoc", *map(str, arguments), "--fail-if-warnings"],
        input=input_text, capture_output=True, text=True, encoding="utf-8",
    )
    require(result.returncode == 0, result.stderr or "Pandoc failed")
    return result.stdout


def text_of(node):
    if isinstance(node, list):
        return "".join(text_of(child) for child in node)
    if not isinstance(node, dict):
        return ""
    if node.get("t") == "Str":
        return node["c"]
    if node.get("t") in {"Space", "SoftBreak", "LineBreak"}:
        return " "
    return text_of(node.get("c", []))


def link(label, target):
    return {"t": "Link", "c": [["", [], []], [{"t": "Str", "c": label}], [target, ""]]}


def assemble():
    chapters = sorted((ROOT / "chapters").glob("[0-9][0-9]-*.md"))
    require(len(chapters) == 10, "Expected ten chapter files")
    paths = [ROOT / "README.md", *chapters, ROOT / "appendix.md"]
    documents = {}
    destinations = {paths[0].resolve(): "contents", paths[-1].resolve(): "appendix"}
    destinations.update((path.resolve(), f"chapter-{index:02d}")
                        for index, path in enumerate(chapters, 1))
    for path in paths:
        text = path.read_text(encoding="utf-8-sig")
        require(not EDITORIAL_MARKERS.search(text), f"Editorial marker in {path.name}")
        if path in chapters:
            require(re.search(r"^## Українськ", text, re.MULTILINE),
                    f"Missing Ukrainian example: {path.name}")
        for match in re.finditer(r"\]\(([^)]+\.md)(?:#[^)]*)?\)", text):
            require((path.parent / unquote(match.group(1))).is_file(),
                    f"Broken manuscript link: {path.name}: {match.group(1)}")
        documents[path.resolve()] = json.loads(pandoc(path, "--from=gfm", "--to=json"))
    references = set()
    heading_targets = {}
    chapter_labels = []
    for path in paths[1:]:
        chapter_headers = 0
        for block in documents[path.resolve()]["blocks"]:
            if block["t"] != "Header":
                continue
            level, attributes, inlines = block["c"]
            original_id = attributes[0]
            label = text_of(inlines)
            reference = re.match(r"([HNAPVDMSRU]\d{2})\.", label)
            if path == paths[-1] and reference:
                require(reference.group(1) not in references, "Duplicate source code")
                references.add(reference.group(1))
                identifier = "ref-" + reference.group(1).lower()
            elif level == 1:
                identifier = destinations[path.resolve()]
                chapter_headers += 1
                if path in chapters:
                    chapter_labels.append((label, identifier))
            else:
                identifier = destinations[path.resolve()] + "-" + original_id
            heading_targets[(path.resolve(), original_id)] = identifier
            attributes[0] = identifier
        require(chapter_headers == 1, f"Expected one title in {path.name}")
    require(len(references) >= 40, "Bibliography is unexpectedly short")

    def rewrite(node, path):
        if isinstance(node, list):
            return [rewrite(child, path) for child in node]
        if not isinstance(node, dict):
            return node
        if node.get("t") == "Header":
            return node
        if node.get("t") == "Link":
            target = urlsplit(node["c"][2][0])
            if not target.scheme and target.path.lower().endswith(".md"):
                destination = (path.parent / unquote(target.path)).resolve()
                require(destination in destinations, f"Link outside edition: {target.path}")
                anchor = destinations[destination]
                if target.fragment:
                    anchor = heading_targets[(destination, unquote(target.fragment))]
                node["c"][2][0] = "#" + anchor
            return node
        if node.get("t") == "Str":
            content = node["c"]
            matches = list(SOURCE_CODE.finditer(content))
            if not matches:
                return node
            pieces = []
            offset = 0
            for match in matches:
                code = match.group()
                require(code in references, f"Unknown citation {code}: {path.name}")
                if offset < match.start():
                    pieces.append({"t": "Str", "c": content[offset:match.start()]})
                pieces.append(link(code, "#ref-" + code.lower()))
                offset = match.end()
            if offset < len(content):
                pieces.append({"t": "Str", "c": content[offset:]})
            return {"t": "Span", "c": [["", [], []], pieces]}
        if "c" in node:
            node["c"] = rewrite(node["c"], path)
        return node

    front = []
    active = False
    first = documents[paths[0].resolve()]
    for block in first["blocks"]:
        if block["t"] == "Header":
            label = text_of(block["c"][2])
            if label == "Від автора":
                active = True
                block["c"][0] = 1
                block["c"][1][0] = "preface"
            elif label == "Зміст":
                break
        if active:
            front.append(block)
    require(front, "Author preface missing")
    items = [("Від автора", "preface"), *chapter_labels,
             ("Словник, абревіатури та джерела", "appendix")]
    contents = {"t": "Div", "c": [
        ["contents", ["book-contents"], [["role", "doc-toc"]]],
        [{"t": "Header", "c": [2, ["contents-heading", [], []], [{"t": "Str", "c": "Зміст"}]]},
         {"t": "BulletList", "c": [[{"t": "Plain", "c": [link(label, "#" + anchor)]}]
                                     for label, anchor in items]}],
    ]}
    first["blocks"] = [contents, *front]
    for path in paths[1:]:
        for block in documents[path.resolve()]["blocks"]:
            if block["t"] == "Para" and any(
                child.get("t") == "Link" and child["c"][2][0].endswith("README.md")
                for child in block["c"] if isinstance(child, dict)
            ):
                block = {"t": "Div", "c": [["", ["chapter-navigation"], []], [block]]}
            first["blocks"].append(rewrite(block, path.resolve()))
    first["meta"] = {key: {"t": "MetaString", "c": value} for key, value in {
        "title": TITLE, "author": AUTHOR, "lang": "uk", "date": "2026-10",
        "subtitle": "Як інформаційні системи вчаться не довіряти чужим командам",
        "toc-title": "Зміст",
    }.items()}
    return first


class HTMLIndex(HTMLParser):
    def __init__(self):
        super().__init__()
        self.ids = set()
        self.links = []
        self.assets = []
        self.language = None

    def handle_starttag(self, tag, attributes):
        attributes = dict(attributes)
        if identifier := attributes.get("id"):
            require(identifier not in self.ids, f"Duplicate HTML id: {identifier}")
            self.ids.add(identifier)
        if tag == "html":
            self.language = attributes.get("lang")
        if tag == "a":
            self.links.append(attributes.get("href", ""))
        if tag in {"script", "img", "link"}:
            resource = attributes.get("src") or attributes.get("href", "")
            if resource and not resource.startswith("data:"):
                self.assets.append(resource)


def check_chapters_and_sources(ids, expected):
    require(all(f"chapter-{index:02d}" in ids for index in range(1, 11)), "Missing chapter")
    require(sum(identifier.startswith("ref-") for identifier in ids) == expected, "Missing sources")
    require("contents" in ids and "appendix" in ids, "Missing contents or appendix")


def validate_html(path, expected):
    parser = HTMLIndex()
    parser.feed(path.read_text(encoding="utf-8"))
    require(parser.language == "uk", "Incorrect HTML language")
    require(not parser.assets, f"HTML requires external resources: {parser.assets}")
    for reference in parser.links:
        target = urlsplit(reference)
        if not target.scheme and not target.netloc:
            require(not target.path, f"Unresolved HTML file link: {reference}")
            require(unquote(target.fragment) in parser.ids, f"Missing HTML anchor: {reference}")
    check_chapters_and_sources(parser.ids, expected)
    print(f"PASS HTML: Ukrainian language, ten chapters, {expected} sources, contents, links, embedded CSS")


def validate_epub(path, expected):
    namespaces = {"container": "urn:oasis:names:tc:opendocument:xmlns:container",
                  "opf": "http://www.idpf.org/2007/opf",
                  "dc": "http://purl.org/dc/elements/1.1/"}
    with zipfile.ZipFile(path) as archive:
        require(archive.testzip() is None, "Corrupt EPUB ZIP entry")
        require(archive.namelist()[0] == "mimetype", "EPUB mimetype must be first")
        require(archive.getinfo("mimetype").compress_type == zipfile.ZIP_STORED,
                "EPUB mimetype must be uncompressed")
        require(archive.read("mimetype") == b"application/epub+zip", "Incorrect EPUB mimetype")
        container = ET.fromstring(archive.read("META-INF/container.xml"))
        package_path = container.find(".//container:rootfile", namespaces).get("full-path")
        package = ET.fromstring(archive.read(package_path))
        require(package.get("version") == "3.0", "Expected EPUB 3")
        for key, value in {"title": TITLE, "creator": AUTHOR, "language": "uk"}.items():
            require(package.findtext(f"opf:metadata/dc:{key}", namespaces=namespaces) == value,
                    f"Incorrect EPUB {key}")
        entries = {}
        documents = {}
        for item in package.findall("opf:manifest/opf:item", namespaces):
            resource = posixpath.normpath(posixpath.join(posixpath.dirname(package_path), unquote(item.get("href"))))
            require(resource in archive.namelist(), f"Missing EPUB manifest resource: {resource}")
            entries[item.get("id")] = (resource, item)
            if item.get("media-type") == "application/xhtml+xml":
                documents[resource] = ET.fromstring(archive.read(resource))
        require(any("nav" in item.get("properties", "").split() for _, item in entries.values()),
                "Missing EPUB navigation")
        for item in package.findall("opf:spine/opf:itemref", namespaces):
            require(item.get("idref") in entries, "Missing EPUB spine item")
        identifiers = {}
        for resource, document in documents.items():
            ids = [element.get("id") for element in document.iter() if element.get("id")]
            require(len(ids) == len(set(ids)), f"Duplicate EPUB id: {resource}")
            identifiers[resource] = set(ids)
        for resource, document in documents.items():
            for element in document.iter():
                for attribute in ("href", "src"):
                    reference = element.get(attribute)
                    if not reference:
                        continue
                    target = urlsplit(reference)
                    if target.scheme or target.netloc:
                        continue
                    destination = posixpath.normpath(posixpath.join(posixpath.dirname(resource), unquote(target.path))) if target.path else resource
                    require(destination in archive.namelist(), f"Broken EPUB resource link: {reference}")
                    if target.fragment:
                        require(unquote(target.fragment) in identifiers.get(destination, set()),
                                f"Broken EPUB fragment: {resource}: {reference}")
        check_chapters_and_sources(set().union(*identifiers.values()), expected)
        print(f"PASS EPUB 3: {len(documents)} XHTML files, ZIP, metadata, manifest, spine, navigation, links")


def main():
    require(shutil.which("pandoc"), "Install Pandoc and add it to PATH")
    serialized = json.dumps(assemble(), ensure_ascii=False)
    appendix = (ROOT / "appendix.md").read_text(encoding="utf-8-sig")
    expected = len(re.findall(r"^### [HNAPVDMSRU]\d{2}\.", appendix, re.MULTILINE))
    destination = ROOT / "dist"
    destination.mkdir(exist_ok=True)
    common = ["--from=json", "--standalone", "--css", ROOT / "book.css"]
    html_path = destination / f"{STEM}.html"
    epub_path = destination / f"{STEM}.epub"
    pandoc(*common, "--to=html5", "--embed-resources", "--section-divs", "--output", html_path, input_text=serialized)
    pandoc(*common, "--to=epub3", "--toc", "--toc-depth=1", "--split-level=1", "--output", epub_path, input_text=serialized)
    validate_html(html_path, expected)
    validate_epub(epub_path, expected)
    plain = pandoc("--from=json", "--to=plain", input_text=serialized)
    words = len(re.findall(r"\b[^\W\d_]+(?:['’-][^\W\d_]+)*\b", plain))
    print(f"Reading edition: approximately {words:,} words, ten chapters and appendix")
    print(html_path)
    print(epub_path)


if __name__ == "__main__":
    main()