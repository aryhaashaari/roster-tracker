#! /bin/sh
echo -n "Enter new go package name (i.e. gitlab.privy.id/privypass/repo-name): "
read pkg_name

if [ -z "$pkg_name" ]; then
    echo "Error: package name cannot be empty"
    exit 1
fi

repo_name=$(echo "$pkg_name" | rev | cut -d'/' -f1 | rev)

# update go mod
echo "update go.mod"
go mod edit -module $pkg_name

echo "converting all imports in code"
escaped_pkg_name=$(echo $pkg_name | sed 's/\//\\\//g')
old_pkg_name=$(echo 'gitlab.privy.id/privypass/privypass-boilerplate' | sed 's/\//\\\//g')
# fix imports project wide, replaces .../go-start with escaped pkg_name
find . -name "*.go" -print0 | xargs -0 sed -i -s "s/${old_pkg_name}/${escaped_pkg_name}/g"

find . -name "README.md.ori" -print0 | xargs -0 sed -i -s "s/privypass-boilerplate/${repo_name}/g"
rm "README.md"
mv "README.md.ori" "README.md"

echo "deleting backup file"
find . -type f -name "*-s" -delete

echo "changing git origin"
git remote remove origin
git remote add origin "https://$pkg_name"

mv "../privypass-boilerplate" "../$repo_name"