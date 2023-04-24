go mod edit -module {'shop'}
find . -type f -name '*.go' \
  -exec sed -i -e 's,{OLD_MODULE},{NEW_MODULE},g' {} \;
