
# maybe more powerful
# for mac (sed for linux is different)
grep "x-real-control" * -R | grep -v Godeps | awk -F: '{print $1}' | sort | uniq | xargs sed -i '' 's#x-real-control#real-liebian#g'
