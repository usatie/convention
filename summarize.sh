#!/bin/bash
cd results
# count
echo "合計:"
echo -e "$(grep -v '^#' results.txt | wc -l | awk '{print $1}') / $(grep '^# ' results.txt | wc -l | awk '{print $1}')\n"

# result
echo "具体的なコード例:"
head -100 results.txt | sed 's/.*: //g' | uniq
