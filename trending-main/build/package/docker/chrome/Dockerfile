FROM ubuntu:jammy

# Add a user for running applications.
RUN useradd apps \
 && mkdir -p /home/apps && chown apps:apps /home/apps \
 && apt-get update && apt-get install -y \
            x11vnc \
            xvfb  \
            fluxbox \
            wget \
            gnupg \
            wmctrl \
            locales \
 && wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - \
 && echo "deb http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google.list \
 && apt-get update && apt-get -y install google-chrome-stable socat \
 && apt-get install -y socat \
 && sed -i '/ru_RU.UTF-8/s/^# //g' /etc/locale.gen \
 && locale-gen \
 && rm -rf /var/lib/apt/lists/*

ENV LANG ru_RU.UTF-8
ENV LANGUAGE ru_RU:ru
ENV LC_ALL ru_RU.UTF-8

COPY build/package/docker/chrome/bootstrap.sh /

EXPOSE 34000

CMD '/bootstrap.sh'
