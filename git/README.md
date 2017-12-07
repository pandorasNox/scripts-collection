delete branches under a certain namespace:
`git branch -r | grep "origin/hotfix/XYZ/" | sed 's/origin\///' | xargs -I %s bash -c "git branch -rd origin/%s; git push origin :%s"`
