FROM sharelatex/sharelatex:latest

RUN apt-get update && \
    echo "ttf-mscorefonts-installer msttcorefonts/accepted-mscorefonts-eula select true" | debconf-set-selections && \
    apt-get install -y ttf-mscorefonts-installer fontconfig && \
    fc-cache -f -v

RUN tlmgr update --self && \
    tlmgr install scheme-full