# locuste.service.osm
LOCUSTE : Service ordonnanceur / Pilotage automatique / Gestionnaire de vol

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/0f2fd8eb4b2149ae85807192e515e7ac)](https://www.codacy.com/manual/axel.maciejewski/locuste.service.osm?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=DaemonToolz/locuste.service.osm&amp;utm_campaign=Badge_Grade)


[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=DaemonToolz_locuste.service.osm&metric=alert_status)](https://sonarcloud.io/dashboard?id=DaemonToolz_locuste.service.osm)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=DaemonToolz_locuste.service.osm&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=DaemonToolz_locuste.service.osm)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=DaemonToolz_locuste.service.osm&metric=reliability_rating)](https://sonarcloud.io/dashboard?id=DaemonToolz_locuste.service.osm)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=DaemonToolz_locuste.service.osm&metric=security_rating)](https://sonarcloud.io/dashboard?id=DaemonToolz_locuste.service.osm)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=DaemonToolz_locuste.service.osm&metric=vulnerabilities)](https://sonarcloud.io/dashboard?id=DaemonToolz_locuste.service.osm)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=DaemonToolz_locuste.service.osm&metric=bugs)](https://sonarcloud.io/dashboard?id=DaemonToolz_locuste.service.osm)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=DaemonToolz_locuste.service.osm&metric=coverage)](https://sonarcloud.io/dashboard?id=DaemonToolz_locuste.service.osm)

Le project Locuste se divise en 3 grandes sections : 
* Automate (Drone Automata) PYTHON
* Unité de contrôle (Brain) GOLANG
* Unité de planification de vol / Ordonanceur (Scheduler) GOLANG
* Interface graphique (UI) ANGULAR


![image](https://user-images.githubusercontent.com/6602774/82243830-8960ca80-9940-11ea-917e-15585f178c6d.png)

Tout le système est embarqué sur une carte Raspberry PI 4B+, Raspbian BUSTER.
* Golang 1.11.2
* Angular 9
* Python 3.7
* Dépendance forte avec la SDK OLYMPE PARROT : (https://developer.parrot.com/docs/olympe/, https://github.com/Parrot-Developers/olympe)


![image](https://user-images.githubusercontent.com/6602774/82240232-59162d80-993a-11ea-8f8e-c7d3cfde2a7c.png)


Détail des choix techniques pour la partie Ordonaceur :

* [Golang] - Conservation de la continuité des développements entrepris par la section [Unité de contrôle]. Il était plus simple, plus rapide de conserver le même langage entre ces deux modules fortement liés par le biais des composants RPC.
* [RPC] - Une des méthodes de communication les plus rapide 

