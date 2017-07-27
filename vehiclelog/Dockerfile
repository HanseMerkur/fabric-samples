FROM node:6.11.1

ENV APPDIR=/app/vehiclelog

# Alle Dateien in den Container kopieren
COPY ./ $APPDIR/
WORKDIR $APPDIR

RUN npm install

RUN chown -R node:node $APPDIR

USER node

EXPOSE 3000
CMD ["npm", "start"]
