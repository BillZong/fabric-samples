CHANNEL_NAME=$1
CONFIG_PATH=$2

# import utils
. scripts/utils.sh

echo "Installing jq"
apt-get -y install jq

pwd

# Fetch the config for the channel, writing it to the path
fetchChannelConfig ${CHANNEL_NAME} ${CONFIG_PATH}
chmod 666 ${CONFIG_PATH}

exit 0
