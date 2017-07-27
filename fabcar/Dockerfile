FROM node:6.11.1

ENV APPDIR=/app/fabcar

# Alle Dateien in den Container kopieren
COPY ./ $APPDIR/
WORKDIR $APPDIR

RUN npm install

RUN chown -R node:node $APPDIR

USER node

CMD ["bash"]
