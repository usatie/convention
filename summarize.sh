#!/bin/bash
cd results

# Examples
echo "具体的なコード例:"
tail -100 results.txt | sed 's/.*: //g' | uniq
echo

# Summary
echo "検出された行数      : $(grep -v '^#' results.txt | wc -l | awk '{print $1}')"
echo "検査したモジュール数: $(grep '^# ' results.txt | wc -l | awk '{print $1}')"
