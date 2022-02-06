Hallo Peer-Reviewende,

In dieser Text-Datei stellen wir unser Projekt erstmal grundlegend vor.
Wir haben beschlossen unsere Website "The Adversary" zu nennen. "The Adversary" ist Englisch
und lässt sich mit "Der Widersacher" übersetzen. Auf unserer Website planen wir Artikel über Themen
die uns interessieren zu veröffentlichen, wie unter anderem Open-Source Projekte, Programmieren
und weitere Sachen die uns interessieren. Wir haben uns jedoch keinen festen Themenrahmen gesetzt.

Wir haben für unsere Webseite ein Front- sowie Backend, das sich in diesem Ordner befindet.

Falls du unter Windows bist haben wir dir das Backend in kompilierter Form beigelegt, weil wir
die mit Compilation einhergehenden Probleme ersparen wollten. Jedoch kannst du auch das Backend selber 
kompilieren. Folgende Programme werden für das kompilieren benötigt:

go  (Programmiersprache mit der das Backend programmiert wurde)
gcc (Einige Go Module haben benötigen C Code der kompiliert werden muss)

Wir haben unsere Website auch schon gehostet. 
Diese ist unter https://bytedream.org/theadversary/ erreichbar.
Um die Website lokal aufzurufen muss die Datei TheAdversary.exe (auf Windows) oder TheAdversary (auf Linux) gestartet
und im Browser http://localhost:8080 aufgerufen werden.

Falls du Hilfe benötigst kannst du dich jederzeit an

- Discord
    demondave#1483 (Ansprechpartner Frontend)
    ByteDream#4312 (Ansprechpartner Backend)

- E-Mail
    davidmaul5@gmail.com (Ansprechpartner Frontend)
    bytedream@protonmail.com (Ansprechpartner Backend)

wenden.

Websitenstruktur:

Wenn du auf die Website gehst siehst du eine Liste mit Artikeln, momentan gibt es nur wenige Artikel,
jedoch haben wir einen Artikel grob vorbereitet. Dieser Artikel sollte "Beispielartikel" heißen,
oben in dem Artikel ist ein Canvas eingebettet, das die Mandelbrot-Menge (https://de.wikipedia.org/wiki/Mandelbrot-Menge)
zeigt diese wird in Echtzeit im Browser mithilfe von WebAssembly (https://de.wikipedia.org/wiki/WebAssembly) 
gerendert. Der dahinterliegende Code ist in Rust (https://de.wikipedia.org/wiki/Rust_(Programmiersprache))
geschrieben. Des weiteren entält unsere Website ein Suchfeld um Artikel zu suchen.

Datei- / und Ordnerstruktur:

MandelbrotWASM:
    MandelbrotWASM enthält den Source Code für den Mandelbrot Renderer.

backend:
    Der backend Ordner enthält den Source Code für das Backend

frontend:
    Der frontend Ordner enhält HTML, CSS, etc. 
    Der /frontend/html/ Ordner enthält die HTML-Dateien.
    Dateien die mit .gohtml enden sind Templates. 
    Die Dateien about.html, contact.html und legal-notice.html sind statische Dateien.

.env:
    Die Informationen wie der Server laufen soll stehen hier drin.
    PORT -> Port auf dem der Server laufen soll.
    ADDRESS -> Adresse über die der Server aufgerufen werden soll.
    SUBPATH -> (Optionaler) Pfad unter dem Server aufgerufen werden soll.
    DATABASE_FILE -> Datenbank in der die Artikel, Autoren und Tags gespeichert werden.
    FRONTEND_DIR -> Pfad zum Ordner der das Frontend enthält.

database.sqlite3:
    Die Datenbank die Artikel, Autoren und Tags gespeichert hat.

Dockerfile:
    Eine Dockerfile (https://de.wikipedia.org/wiki/Docker_(Software)) mit der der Server als Docker
    Container aufgesetzt werden kann.
    Zum Starten auf Linux kann der Folgende Befehl verwendet werden:
        >>> docker build -t theadversary . && docker run -p 8080:8080 theadversary

TheAdversary (Linux) / TheAdversary (Windows):
    Das kompilierte Backend aus dem backend/ Ordner.


Bei weiteren Fragen können die Ansprechpartner der jeweiligen Teams über die oben angegebenen
Kontaktmöglichkeiten gerne kontaktiert werden.
