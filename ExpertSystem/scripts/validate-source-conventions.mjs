#!/usr/bin/env node

import { readdir, readFile } from "node:fs/promises";
import process from "node:process";

const root = process.cwd();
const failures = [];

function report(file, line, message) {
  failures.push(`${file}:${line} ${message}`);
}

async function main() {
  const files = (await readdir(root))
    .filter((name) => name.toLowerCase().endsWith(".md"))
    .sort();

  for (const file of files) {
    const source = await readFile(new URL(`../${file}`, import.meta.url), "utf8");
    const lines = source.split(/\r?\n/u);

    for (const [index, line] of lines.entries()) {
      const lineNumber = index + 1;

      if (/^\s*\$\$\s*$/u.test(line)) {
        report(file, lineNumber, "старий дисплейний роздільник $$; використайте fenced-блок ```math");
      }

      if (/\\operatorname\b/u.test(line)) {
        report(file, lineNumber, "макрос \\operatorname несумісний із цільовим рендерером; використайте \\mathrm");
      }

      if (/екосистем/iu.test(line)) {
        report(file, lineNumber, "заборонене слово «екосистема» або його словоформа");
      }

      if (/Питання до читачів\s*[:.]/u.test(line)) {
        report(file, lineNumber, "оформіть «Питання до читачів» заголовком без двокрапки або крапки");
      }
    }
  }

  if (failures.length > 0) {
    console.error(`Source convention validation: ${failures.length} error(s)`);
    for (const failure of failures) console.error(`- ${failure}`);
    process.exitCode = 1;
    return;
  }

  console.log(`Source convention validation passed: ${files.length} files.`);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
