## Digitales Fahrtenbuch -- Vehicle Log

Diese Beispielanwendung umfasst eine Node.js Anwendung und demonstriert den Umgang mit den **_fabric-client_** & **_fabric-ca-client_** Node.js SDK APIs. Die Anwendung ist so konzipiert, dass sie komplett unter Windows mit der Docker Toolbox ausgeführt wird. Es ist keine lokale Installation von **_Go_** oder **_Node.js_** notwendig.

### Voraussetzungen:

* [Docker](https://www.docker.com/products/overview) - v1.12 or höher
* [Docker Compose](https://docs.docker.com/compose/overview/) - v1.8 or höher
* [Git client](https://git-scm.com/downloads) - wird für den clone Befehl benötigt
* [Docker Toolbox] (https://docs.docker.com/toolbox/overview) - enthält Docker und Docker Compose für ältere Windows System (Windows 7 oder Windows 10 Home)

### Installation
Die Installation von Docker und Docker Compose für Linux Umgebungen ist auf der Homepage von Docker sehr gut beschrieben, z.B. unter [Get Docker CE for Debian](https://docs.docker.com/engine/installation/linux/docker-ce/debian/).

Für Windows-Umgebungen und insbesondere für alte Windows-Umgebungen oder Windows 10 Home kann Docker nicht direkt ausgeführt werden, sondern es wird eine eigene Virtuelle Maschine benötigt, die mit Oracle Virtual Box ausgeführt wird. Dies Umgebung heißt Docker Toolbox.

Bei der Verwendung mit der Docker Toolbox müssen einige Dinge beachtet werden, damit alle Aufrufe erfolgreich durchgeführt werden können. Im Folgenden wird daher beschrieben, wie initial eine Umgebung aufgebaut wird.

#### Eigene Docker Umegbung für Hyperledger
Nach der Installation von Docker Toolbox existiert bereits eine Virtuelle Maschine mit dem Namen ```default```. Diese Maschine bleibt unverändert, statt dessen wird eine zweite Docker Umgebung mit dem Namen ```hyperledger``` eingerichtet. Dazu müssen die folgenden Befehle in einem Git Bash Terminal ausgeführt werden:

```
docker-machine.exe create --driver virtualbox hyperledger
```

#### Shared Folder außerhalb des Benutzerprofils
Für diese neue Umgebung wird automatisch ein Shared Folder für das Benutzer-Verzeichnis unter Windows (```C:\Users```) eingerichtet. Falls eigenen Projekte nicht im Benutzerprofil liegen sollen, wie es in Unternehmensumgebungen oft der Fall ist, dann muss ein weiterer Shared Folder eingerichtet werden. Dieser Schritt ist sehr wichtig, da andernfalls die Einbindung von Volumes unter Docker nicht funktioniert. 
Der folgende Code zeigt die Befehle die notwendig sind, um den Ordner ```C:\anw_entw``` als Shared Folder verfügbar zu machen und in der Docker Umgebung unter dem Pfad ```/c/anw_entw``` verfügbar zu machen:

```
docker-machine.exe stop hyperledger
/c/Program\ Files/Oracle/VirtualBox/VBoxManage.exe sharedfolder add
hyperledger --name "anw_entw" --hostpath "C:\anw_entw" --automount
docker-machine.exe start hyperledger
docker-machine.exe ssh hyperledger "sudo sh -c 'echo \"#!/bin/sh\" >
/var/lib/boot2docker/bootlocal.sh'"
docker-machine.exe ssh hyperledger "sudo sh -c 'echo \"mkdir -p /c/anw_entw\" >>
/var/lib/boot2docker/bootlocal.sh'"
docker-machine.exe ssh hyperledger "sudo sh -c 'echo \"mount -t vboxsf -o
uid=1000,gid=50 anw_entw /c/anw_entw\" >> /var/lib/boot2docker/bootlocal.sh'"
docker-machine.exe ssh hyperledger "sudo chmod u+x
/var/lib/boot2docker/bootlocal.sh"
docker-machine.exe ssh hyperledger "sudo /var/lib/boot2docker/bootlocal.sh"
```
Über die Anwendung ```VBoxManage``` wird der Ordner ```C:\anw_entw``` als Shared Folder eingerichtet. Im Anschluss wird innerhalb der Docker Umgebung eine Datei ```/var/lib/boot2docker/bootlocal.sh``` angelegt. In dieser Datei wird dafür gesorgt, dass beim Start der Docker Umgebung der Shared Folder unter dem Pfad ```/c/anw_entw``` gemountet wird. Hierbei ist es ganz wichtig, dass der Pfad genauso lautet, wie er auch im Git Bash Terminal lautet.

#### Proxy einrichten
Falls ein Internet-Proxy den Zugang zum Internet schützt, so muss auch die Docker Umgebung diesen Proxy kennen. Die folgenden beiden Codezeilen fügen den Proxy zu einer Konfigurationsdatei hinzu

```
docker-machine.exe ssh hyperledger "sudo sh -c 'echo \"export HTTP_PROXY=http://proxy.example.com:8080\" >> /var/lib/boot2docker/profile'"
docker-machine.exe ssh hyperledger "sudo sh -c 'echo \"export HTTPS_PROXY=http://proxy.example.com:8080\" >> /var/lib/boot2docker/profile'"
```

#### Docker Umgebung testen
Ein einfacher Test der Docker Umgebung sieht wie folgt aus:

```
docker-machine start hyperledger
eval $(docker-machine env hyperledger)
docker run --rm hello-world
```

### Ausführung des Beispiels

#### Chaincode übersetzen
Eine ausführliche Beschreibung zum Thema Chaincode findet sich unter [Chaincode for Developers] (https://hyperledger-fabric.readthedocs.io/en/latest/chaincode4ade.html).

Der Chaincode für das Beispiel ist in der Datei ```fabric-samples/chaincode/vehiclelog/vehiclelog.go``` enthalten. Nach Änderungen an dieser Datei ist es sinnvoll, sie einmal zu übersetzen, um Syntaxfehler zu vermeiden. Auch dies passiert in einem Docker-Container, so dass keine Installation von ***Go*** notwendig ist.

Für die Übersetzung muss ein *Git Bash Terminal* geöffnet werden und dort die folgenden Befehle ausgeführt werden:
```
eval $(docker-machine.exe env hyperledger)
cd /anw_entw/blockchain/fabric-samples/chaincode
docker-compose run --rm chaincode bash -c 'cd vehiclelog && go build -x'
```
Der erste Befehl ist nur für die Ausführung mit Docker Toolbox notwendig, um die Umgebungsvariablen zu setzen. Nach dem Wechsel in das passende Verzeichnis wird über ```docker-compose``` ein kurzlebiger One-Off Container mit dem ```go build``` Befehl ausgeführt. Kann der Chaincode erfolgreich übersetzt werden, so wird eine ausführbare Datei ```vehiclelog``` angelegt. Diese wird später nicht mehr benötigt und kann gelöscht werden.

#### Blockchain Netzwerk starten
Für die Ausführung der Beispiel-Blockchain bzw. des zugehörigen Chaincodes wird ein Blockchain-Netzwerk benötigt. In den Beispielen ist dafür mit dem ```basic-network``` ein einfaches Netzwerk vorkonfiguriert. Es besteht aus insgesamt 4 Docker-Containern. Gestartet wird das Netzwerk wie folgt (das Git Bash Terminal von oben kann benutzt werden):
```
cd /anw_entw/blockchain/fabric-samples/basic-network
./start.sh
```
Das Skript ```start.sh``` beendet zunächst per ```docker-compose down``` die Container, falls sie von einer früheren Sitzung noch laufen. Danach werden über ```docker-compose up``` vier Container gestartet:
* ein ```ca.example.com``` Container mit einer Certificate Authority für die Überprüfung von Zertifikaten,
* ein ```orderer.example.com``` Container für die Bearbeitung von Aufträgen,
* ein ```peer0.org1.example.com``` Container als Client-Peer der Blockchain,
* ein ```couchdb.example.com```Container für die Speicherung der Blockchain-Daten
Im Anschluss erzeugt das Skript einen neuen Channel ```mychannel``` und fügt den Client-Peer diesem Channel hinzu.

#### Chaincode in den Channel deployen
Der oben übersetzte Chaincode muss vor der Ausführung in das Blockchain Netzwerk deployed werden. Dies passiert über einen eigenen ```cli```Docker Container, der passend zum ```basic-network``` konfiguriert ist. Das Deployment besteht aus den drei einzelnen Schritten Installation, Instanziierung und Initialisierung.
```
cd /anw_entw/blockchain/fabric-samples/basic-network
docker-compose -f ./docker-compose.yml up -d cli

docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=//opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp" cli peer chaincode install -n vehiclelog -v 1.0 -p github.com/vehiclelog

docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=//opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp" cli peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n vehiclelog -v 1.0 -c '{"Args":[""]}' -P "OR ('Org1MSP.member','Org2MSP.member')"

sleep 10

docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=//opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp" cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n vehiclelog -c '{"function":"initLedger","Args":[""]}'
```
Für die Befehle aus den beiden vorherigen Abschnitten *Netzwerk starten* sowie *Chaincode deployen* gibt es im Verzeichnis ```.../fabric-samples/vehiclelog``` auch das Skript ```startFabric.sh````, welches die Einzelschritte nacheinander ausführt.

#### Chaincode testen
Sobald der Chaincode erfolgreich deployed wurde, ist es Zeit für einen ersten Test. Dazu kann die kleine Node.js Anwendung ```query.js``` benutzt werden, die eine ganz einfache Anfrage an die Blockchain mit dem Chaincode sendet:
```
cd /c/anw_entw/blockchain/fabric-samples/vehiclelog
docker-compose build
docker-compose run --rm vehiclelog node query.js
```
Mit diesen Befehlen wird zunächst ein einfaches Docker Image auf Basis von Node.js gebaut und mit den notwendigen npm-Packages gefüllt (der *build* Aufruf). Mit diesem Image wird dann ein Container gestartet, der die Datei ```query.js``` ausführt.

Der Aufruf liefert eine Liste aller Fahrzeuge, die in der Blockchain hinterlegt sind. Dies sieht z.B. wie folgt aus:

```
Create a client and set the wallet location
Set wallet path, and associate user  PeerAdmin  with application
Check user is enrolled, and set a query URL in the network
Make query
Assigning transaction_id:  f790cd76e99453228a504b4cb6f115a9cf30ebccd02e3238d8eb26942e50734d
returned from query
Query result count =  1
Response is  [{"Key":"FIN-001", "Record":{"Registration":"ZWS-IT 001","manufactor":"VW","model":"California","type":"Kfz"}},{"Key":"FIN-002", "Record":{"Registration":"ZWS-IT 002","manufactor":"BMW","model":"R 1100 R","type":"Krad"}},{"Key":"FIN-003", "Record":{"Registration":"ZWS-IT 003","manufactor":"Honda","model":"NTV","type":"Krad"}}]
```

##### Hinweis
Die npm-Packages werden in das Verzeichnis ```node_modules``` installiert. Dies kann unter Windows zu Problemen führen, weil durch geschachtelte Packages zu lange Pfadnamen entstehen. Dies wird hier dadurch gelöst, dass das ```node_modules```-Verzeichnis in der Datei ```docker-compose.yml``` als Volume Container gemountet wird. Ein ähnliches Vorgehen ist notwendig, um das Verzeichnis ```creds``` unter Windows schreibbar zur Verfügung zu haben.

#### Die Anwendung starten und aufrufen
Eine Node.js Anwendung wird verwendet, um aufbauend auf dem Blockchain Netzwerk und dem Chaincode mit dem Browser auf die Blockcahin zugreifen zu können. Die Node.js Anwendung wird mit dem Docker Image ausgeführt, welches auch für den Test verwendet wurde. Der Start erfolgt wie folgt:

```
cd /c/anw_entw/blockchain/fabric-samples/vehiclelog
docker-compose up -d
```

Das Flag ```-d``` sorgt dafür, dass der Container im Hintergrund weiterläuft. Um die Log-Ausgaben zu sehen. kann der befehl ```docker-compose logs```verwendet werden.
Die Web-Anwendung ist im Anschluss über den Port 8080 des Rechners, auf dem Docker läuft, erreichbar. Dies kann der lokale Rechner sein, also http://localhost:8080 oder, unter Windows, die IP-Adresse der Docker Toolbox VM, also z.B. http://192.168.99.100:8080. 



> Written with [StackEdit](https://stackedit.io/).

