for ARGUMENT in "$@"
do

KEY=$(echo "$ARGUMENT" | cut -f1 -d=)
VALUE=$(echo "$ARGUMENT" | cut -f2 -d=)

case "$KEY" in
REPO)    REPO=${VALUE} ;;
*)
esac
done
PASCAL_REPO=`python3 -c "print('$REPO'.replace('-', ' ').title().replace(' ', ''))"`
CAMEL_REPO=`python3 -c "print('$PASCAL_REPO'[0].lower() + '$PASCAL_REPO'[1:])"`
find .idea/modules.xml -type f -exec sed -i "" -e "s/template-repository-go.iml/$REPO.iml/g" {} +
find .idea/workspace.xml -type f -exec sed -i "" -e "s/template-repository-go.iml/$REPO.iml/g" {} +
find catalog-info.meta.json -type f -exec sed -i "" -e "s/template-repository-go/$REPO/g" {} +
find catalog-info.meta.json -type f -exec sed -i "" -e "s/template/application/g" {} +
find cmd -type f -exec sed -i "" -e "s/honestbank\/template-repository-go/honestbank\/$REPO/g" {} +
find config -type f -exec sed -i "" -e "s/honestbank\/template-repository-go/honestbank\/$REPO/g" {} +
find examples -type f -exec sed -i "" -e "s/honestbank\/template-repository-go/honestbank\/$REPO/g" {} +
find go.mod -type f -exec sed -i "" -e "s/honestbank\/template-repository-go/honestbank\/$REPO/g" {} +
find sonar-project.properties -type f -exec sed -i "" -e "s/honestbank_template-repository-go/honestbank_$REPO/g" {} +
find README.md -type f -exec sed -i "" -e "s/template-repository-go/$REPO/g" {} +
find tools.go -type f -exec sed -i "" -e "s/honestbank\/template-repository-go/honestbank\/$REPO/g" {} +
git mv .idea/template-repository-go.iml .idea/$REPO.iml
rm -Rf .github/workflows/init.yaml
