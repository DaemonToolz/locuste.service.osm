# Locuste Map Scheduler : locuste.service.osm
LOCUSTE : Service ordonnanceur / Pilotage automatique / Gestionnaire de vol

<img width="2609" alt="locuste-scheduler-banner" src="https://user-images.githubusercontent.com/6602774/84285951-5aec9e80-ab3e-11ea-84a9-b5b8dd2f8aed.png">

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/0f2fd8eb4b2149ae85807192e515e7ac)](https://www.codacy.com/manual/axel.maciejewski/locuste.service.osm?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=DaemonToolz/locuste.service.osm&amp;utm_campaign=Badge_Grade)

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=DaemonToolz_locuste.service.osm&metric=alert_status)](https://sonarcloud.io/dashboard?id=DaemonToolz_locuste.service.osm)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=DaemonToolz_locuste.service.osm&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=DaemonToolz_locuste.service.osm)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=DaemonToolz_locuste.service.osm&metric=reliability_rating)](https://sonarcloud.io/dashboard?id=DaemonToolz_locuste.service.osm)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=DaemonToolz_locuste.service.osm&metric=security_rating)](https://sonarcloud.io/dashboard?id=DaemonToolz_locuste.service.osm)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=DaemonToolz_locuste.service.osm&metric=vulnerabilities)](https://sonarcloud.io/dashboard?id=DaemonToolz_locuste.service.osm)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=DaemonToolz_locuste.service.osm&metric=bugs)](https://sonarcloud.io/dashboard?id=DaemonToolz_locuste.service.osm)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=DaemonToolz_locuste.service.osm&metric=coverage)](https://sonarcloud.io/dashboard?id=DaemonToolz_locuste.service.osm)

Le project Locuste se divise en 4 grandes sections : 
* Automate (Drone Automata) PYTHON (https://github.com/DaemonToolz/locuste.drone.automata)
* Unité de contrôle (Brain) GOLANG (https://github.com/DaemonToolz/locuste.service.brain)
* Unité de planification de vol / Ordonanceur (Scheduler) GOLANG (https://github.com/DaemonToolz/locuste.service.osm)
* Interface graphique (UI) ANGULAR (https://github.com/DaemonToolz/locuste.dashboard.ui)

![Composants](https://user-images.githubusercontent.com/6602774/83644711-dcc65000-a5b1-11ea-8661-977931bb6a9c.png)

Tout le système est embarqué sur une carte Raspberry PI 4B+, Raspbian BUSTER.
* Golang 1.11.2
* Angular 9
* Python 3.7
* Dépendance forte avec la SDK OLYMPE PARROT : (https://developer.parrot.com/docs/olympe/, https://github.com/Parrot-Developers/olympe)

![Vue globale](https://user-images.githubusercontent.com/6602774/83644783-f10a4d00-a5b1-11ea-8fed-80c3b76f1b00.png)

Détail des choix techniques pour la partie Ordonanceur :

* [Golang] - Conservation de la continuité des développements entrepris par la section [Unité de contrôle]. Il était plus simple, plus rapide de conserver le même langage entre ces deux modules fortement liés par le biais des composants RPC.
* [RPC] - Une des méthodes de communication les plus rapide 


![Détail de l'initialisation](https://user-images.githubusercontent.com/6602774/82245150-b910d200-9942-11ea-83ab-815dd1db7ee8.png)
