import os
import re
from collections import defaultdict

def generate_link(file, header=None):
    if header:
        header_text = header.strip("#").strip()
        header_id = header_text.lower().replace(" ", "-")
    else:
        header_text = file[:-3]  # Remove file extension
        header_id = ""
    link = f"[{header_text}]({file}#{header_id})"
    return link

def generate_toc(file):
    toc = defaultdict(list)
    main_heading = None
    with open(file, "r") as md_file:
        lines = md_file.readlines()
        for line in lines:
            if line.startswith("#"):
                level = line.count("#", 0, 3)  # only consider levels up to ###
                link = generate_link(file, line)
                if level == 1:
                    main_heading = link
                elif main_heading:  # Only add subheadings if there is a main heading
                    toc[main_heading].append((level, link))

    # Add the file name as a main heading if there is no main heading
    if not main_heading:
        main_heading = generate_link(file)
        toc[main_heading] = []

    return toc

def main():
    toc_filename = "table_of_contents.md"
    toc_dict = defaultdict(list)

    for file in os.listdir("."):
        if file.endswith(".md") and file != toc_filename:
            toc_file = generate_toc(file)
            for main_heading, subheadings in toc_file.items():
                toc_dict[main_heading].extend(subheadings)

    sorted_toc = sorted(toc_dict.items(), key=lambda x: x[0].lower())

    with open(toc_filename, "w") as toc_file:
        toc_file.write("Table of Contents:\n")
        for main_heading, subheadings in sorted_toc:
            toc_file.write(f"- {main_heading}\n")
            for level, link in subheadings:
                indent = "    " * (level - 1)
                toc_file.write(f"{indent}- {link}\n")

if __name__ == "__main__":
    main()
