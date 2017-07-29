FROM node:6.11.1

ENV APPDIR=/app/vehiclelog

# Vorbereitung fuer npm install
COPY package.json $APPDIR/
WORKDIR $APPDIR

RUN npm install

# Die vorhandenen creds-Dateien schreibbar kopieren
COPY ./creds/ $APPDIR/creds/
RUN chown -R node:node $APPDIR

# und auch in den keystore des Benutzers
COPY ./creds/ /home/node/.hfc-key-store/
RUN chown -R node:node /home/node/.hfc-key-store

USER node

EXPOSE 3000
CMD ["npm", "start"]
