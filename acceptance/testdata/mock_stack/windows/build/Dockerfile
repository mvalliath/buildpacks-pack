FROM mcr.microsoft.com/windows/nanoserver:1809

# non-zero sets all user-owned directories to BUILTIN\Users
ENV CNB_USER_ID=1
ENV CNB_GROUP_ID=1

LABEL io.buildpacks.stack.id=pack.test.stack
LABEL io.buildpacks.stack.mixins="[\"mixinA\", \"build:mixinTwo\", \"netcat\", \"mixin3\"]"

USER ContainerUser
