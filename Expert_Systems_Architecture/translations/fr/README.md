# Architecture des systèmes experts fondés sur des preuves

## Résumé

Cet ouvrage étudie comment transformer des documents d’ingénierie, des exigences, du code source et des résultats d’essais en connaissances structurées permettant un raisonnement informatique vérifiable. Réparti en sept parties et quarante chapitres, il traite de la formalisation des concepts et des sources, des modèles et de l’acquisition des connaissances, de leur traçabilité, de l’architecture des systèmes experts, des explications, de la vérification, de la sûreté et des approches neuro-symboliques.

L’ouvrage insiste sur la provenance des éléments probants, leurs limites d’applicabilité et les liens explicites entre exigences, règles, décisions de conception et essais. Il distingue une hypothèse proposée, une conclusion dérivée dans un modèle formel et un résultat confirmé expérimentalement : la seule plausibilité ne constitue pas une preuve. Il aborde également le raisonnement logique, l’incertitude, les modèles causaux et temporels, l’explication des conclusions, les actions contrôlées, l’apprentissage continu et les argumentaires de sûreté. Les modèles de langage sont présentés comme des outils auxiliaires dont les propositions doivent être vérifiées indépendamment.

Le livre s’adresse aux étudiants, ingénieurs et chercheurs qui conçoivent ou évaluent des systèmes d’aide à la décision dans des environnements techniques complexes. Il montre comment créer des systèmes qui conservent les fondements de leurs conclusions, en exposent les limites et savent s’abstenir lorsque les éléments disponibles sont insuffisants.

**Résumés des chapitres :** [les 40 chapitres](#résumés-des-chapitres).

**Mots-clés :** systèmes experts, ingénierie des connaissances, provenance des preuves, traçabilité, vérification formelle, sûreté fonctionnelle, IA neuro-symbolique.

[Édition maîtresse ukrainienne](../../README.md)

---

## Résumés des chapitres

Résumés synthétiques des 40 chapitres de l’édition maîtresse ukrainienne. Chaque numéro renvoie au texte source.

### 01 · [Chapitre 1](../../ch01-introduction-to-expert-systems.md)

Le chapitre définit le système expert par des connaissances explicites, un raisonnement explicable et la capacité de s’abstenir de façon justifiée. Il retrace brièvement son histoire et son architecture, en le distinguant de la recherche documentaire et de la génération libre de texte.

### 02 · [Chapitre 2](../../ch02-epistemology-of-machine-knowledge.md)

Il examine les fondements d’une affirmation produite par une machine : source, contexte, domaine d’application, incertitude et validité. Des limites de température contradictoires montrent qu’un nombre exact ne suffit pas à justifier une décision.

### 03 · [Chapitre 3](../../ch03-beyond-reference-information-systems.md)

Le chapitre distingue systèmes de référence, de recherche, de dialogue, de recommandation et systèmes experts selon leurs tâches et leurs contrats de résultat. Trouver une règle ne signifie pas savoir l’appliquer correctement à un cas précis.

### 04 · [Chapitre 4](../../ch04-evolution-from-bayes-to-evidence-ai.md)

Il retrace l’évolution des systèmes experts et compare logique, probabilités, ensembles flous, graphes et raisonnement par cas. Le choix d’une méthode dépend de la question, des preuves nécessaires et des limites de l’inférence.

### 05 · [Chapitre 5](../../ch05-triad-of-trust-and-corporate-memory.md)

Le chapitre présente un triptyque de confiance : le système applique les connaissances, la recommandation expose ses fondements et la mémoire organisationnelle conserve le contexte et les conséquences. Il traite de provenance, d’audit, d’accès et de révision des connaissances.

### 06 · [Chapitre 6](../../ch06-applied-mathematics-for-expert-systems.md)

À partir de problèmes d’ingénierie, il introduit l’inférence logique et probabiliste, les méthodes par graphes, l’optimisation et l’analyse causale. Un calcul reproductible ou une citation exacte dépend néanmoins de données valides et d’une interprétation correcte.

### 07 · [Chapitre 7](../../ch07-knowledge-base-typology.md)

Le chapitre compare règles, modèles d’objets, ontologies, cas, contraintes, modèles probabilistes et index documentaires ou vectoriels. Il explique ce que chaque représentation conserve et pourquoi les connaissances manquantes doivent rester explicites.

### 08 · [Chapitre 8](../../ch08-engineering-artifacts-as-data.md)

Exigences, décisions, plans et résultats d’essais sont modélisés comme des artefacts structurés, avec des champs typés et des relations. La provenance facilite les requêtes et l’analyse, mais n’établit ni la vérité ni la suffisance juridique d’un document.

### 09 · [Chapitre 9](../../ch09-engineering-knowledge-graph-traceability.md)

Un graphe de connaissances typé permet de suivre les changements, des exigences aux logiciels, au matériel et aux essais. Les requêtes facilitent l’analyse d’impact, mais les liens du graphe ne prouvent pas à eux seuls la conformité.

### 10 · [Chapitre 10](../../ch10-knowledge-acquisition-systems.md)

Le chapitre décrit un processus d’acquisition qui repère les sources, en préserve le contenu et la provenance, vérifie les connaissances candidates et en définit le champ d’application. Fidélité d’une citation, confiance d’un détecteur et approbation documentaire répondent à des questions différentes.

### 11 · [Chapitre 11](../../ch11-knowledge-elicitation-from-experts.md)

Entretiens, comparaison de cas, reconstitution de décisions, cartes cognitives et contre-exemples aident à expliciter le savoir tacite des experts. Fréquence, incertitude d’une estimation et accord entre experts sont traités comme des mesures distinctes.

### 12 · [Chapitre 12](../../ch12-linguistic-analysis-and-local-models.md)

Le chapitre montre comment l’analyse linguistique et les modèles locaux produisent des candidats de règles tout en préservant les octets source, les conditions, la négation et la confidentialité. Une citation exacte ne garantit pas, à elle seule, une interprétation correcte.

### 13 · [Chapitre 13](../../ch13-language-variability-vs-determinism.md)

La recherche de candidats est séparée de la validation sémantique déterministe. Des formulations variées sont converties en requêtes typées, avec vérification des unités, de la négation, de la modalité, de la version des sources et des preuves.

### 14 · [Chapitre 14](../../ch14-requirements-detection-and-formalization.md)

Il traite de la modalité et du champ d’application des exigences, de leur transformation en candidats formels et de la détection des contradictions. Les modèles EARS et les outils limités facilitent la revue sans la remplacer.

### 15 · [Chapitre 15](../../ch15-knowledge-extraction-and-kb-construction.md)

Le chapitre suit le passage de la source au candidat, à l’entrée approuvée de la base de connaissances puis à l’inférence vérifiable. Il décrit l’extraction, les coordonnées d’annotation, les unités et les contrôles d’importation ; les méthodes proposées restent à évaluer.

### 16 · [Chapitre 16](../../ch16-expert-systems-architecture.md)

Il explique les frontières entre dépôts de connaissances, inférence, vérification, politique d’accès et journaux de provenance. Il aborde les mises à jour et les révocations ; les scores de recherche ne sont pas assimilés à des preuves.

### 17 · [Chapitre 17](../../ch17-implementation-stack.md)

Le choix de la pile technologique part des opérations requises, des critères de correction et du comportement en cas d’échec. Langages, moteurs de règles, stockage et déploiement sont comparés, avec une gestion explicite des erreurs.

### 18 · [Chapitre 18](../../ch18-execution-infrastructure.md)

L’architecture est reliée aux contraintes physiques d’exécution : mémoire, processeurs et accélérateurs, bus, latence et emplacement des composants. Les estimations de ressources dépendent du contexte et ne constituent pas des garanties universelles de performance.

### 19 · [Chapitre 19](../../ch19-from-question-to-evidence.md)

Une question devient une affirmation vérifiée en retrouvant ses sources et en la liant à la bonne version, au bon objet, à la bonne unité et à la bonne condition. La pertinence d’un résultat de recherche ne suffit pas à étayer une conclusion.

### 20 · [Chapitre 20](../../ch20-explanation-engine.md)

À partir des traces de décision, le moteur produit des explications répondant à « comment », « pourquoi », « pourquoi pas », « et si » ou « que faut-il changer ». Il traite des conflits minimaux et des accès par rôle, sans laisser un modèle de langage inventer des raisons.

### 21 · [Chapitre 21](../../ch21-from-recommendation-to-action.md)

La recommandation est séparée de l’autorité d’exécuter une action. Une passerelle proposée vérifie l’action, l’acteur, le contexte et la politique, puis journalise l’exécution ; sa conception n’est pas présentée comme validée.

### 22 · [Chapitre 22](../../ch22-cybernetics-edge-to-backend.md)

Un scénario de capteur industriel illustre une rétroaction cybernétique contextualisée, de la mesure à l’analyse en périphérie puis à la réponse. Le score d’un modèle ne suffit pas à décider sans tenir compte de l’étalonnage et de l’état du processus.

### 23 · [Chapitre 23](../../ch23-knowledge-base-verification.md)

Le chapitre distingue la vérification de la base de connaissances implémentée de la validation de la politique au regard de la tâche réelle. Modèles finis, contre-exemples et critères de mise en production révèlent des omissions, sans prouver une correction universelle.

### 24 · [Chapitre 24](../../ch24-system-diagnosis.md)

Le diagnostic fondé sur un modèle compare des hypothèses de panne aux observations lorsque les données sont incomplètes. La cohérence réduit le nombre de candidats, mais ne démontre pas une cause physique sans preuves supplémentaires.

### 25 · [Chapitre 25](../../ch25-how-expert-systems-learn.md)

Matrices d’examen, audits des connaissances, tests de régression et contrôles de calibration soutiennent les mises à jour maîtrisées. Ils détectent les dégradations dans les scénarios testés, sans garantir la correction de tous les cas futurs.

### 26 · [Chapitre 26](../../ch26-continual-learning.md)

Le chapitre explique le biais de rétroaction lorsque les recommandations d’un système influencent les journaux utilisés pour son apprentissage. Suivi et mises à jour contrôlées évitent de confondre choix répétés et preuves indépendantes.

### 27 · [Chapitre 27](../../ch27-safety-case-gsn-synthesis.md)

Des preuves traçables sont associées à une structure d’argumentation de sûreté pour produire un dossier examinable. Il s’agit d’un argument candidat, non d’un certificat automatique ; hypothèses et lacunes doivent être contrôlées.

### 28 · [Chapitre 28](../../ch28-dual-mode-expert-systems.md)

Les conclusions strictes fondées sur des connaissances vérifiées sont séparées des hypothèses consultatives sur les vérifications à mener. Modes, étiquettes et contrôles d’admission empêchent de présenter une suggestion comme une preuve.

### 29 · [Chapitre 29](../../ch29-neuro-symbolic-architecture.md)

Un modèle de langage interprète une question et propose des passages ; un noyau symbolique en vérifie la structure, la provenance, les règles et les exceptions. La preuve d’une conclusion et l’autorisation d’agir restent deux questions distinctes.

### 30 · [Chapitre 30](../../ch30-safety-cybersecurity-co-engineering.md)

La sûreté fonctionnelle et la cybersécurité sont analysées conjointement lorsqu’une mesure de protection peut elle-même créer un risque physique. Contexte partagé et analyse d’impact révèlent les compromis sans subordonner une discipline à l’autre.

### 31 · [Chapitre 31](../../ch31-syllogistic-reasoning-and-relation-lattices.md)

Le chapitre examine le raisonnement normatif par hiérarchies de prédicats, exceptions, réfutations et validité temporelle. Il distingue la recherche d’un passage pertinent de la dérivation d’une conclusion contextuelle fondée sur des règles.

### 32 · [Chapitre 32](../../ch32-high-performance-knowledge-packs-mmap-and-harvesting.md)

Sont présentés des paquets de connaissances immuables, l’admission au niveau des octets, les index et le mappage mémoire, pour un accès reproductible et efficace. Les empreintes protègent l’intégrité, mais ne prouvent ni la vérité ni l’applicabilité du contenu.

### 33 · [Chapitre 33](../../ch33-inter-system-knowledge-exchange-and-model-teaching.md)

Un échange sûr transporte avec chaque valeur les unités, le champ d’application, la version des règles, les sources et les permissions. Les signatures protègent l’origine, pas la correction ni l’applicabilité.

### 34 · [Chapitre 34](../../ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md)

La recherche relationnelle et l’abduction font émerger des explications candidates lorsque des faits manquent, puis des questions ciblées cherchent des éléments discriminants. Les hypothèses restent distinctes des faits établis.

### 35 · [Chapitre 35](../../ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md)

Les systèmes réactifs révisent leurs conclusions lorsque les observations qui les étayaient sont retirées. Les règles immuables sont séparées des données évolutives ; les paramètres adaptatifs restent des propositions à encadrer et à vérifier.

### 36 · [Chapitre 36](../../ch36-knowledge-testing-pyramid-and-variational-calibration.md)

Le concept de pyramide de tests est adapté aux systèmes de connaissances : règles, interactions, parcours complet de réponse et perturbations des entrées. Le chapitre précise ce que chaque niveau permet d’observer et ne permet pas de prouver.

### 37 · [Chapitre 37](../../ch37-input-information-assessment-and-algorithmic-skepticism.md)

Avant l’application des règles, les éléments entrants sont évalués selon l’identité de l’objet, le temps, la provenance, l’indépendance et l’incertitude. Un soupçon déclenche une vérification ; il ne prouve pas une falsification.

### 38 · [Chapitre 38](../../ch38-curing-machine-hallucinations-and-knowledge-deficits.md)

Les propositions d’un modèle de langage sont séparées des conclusions acceptées par des contrôles de source, de sens, d’applicabilité et de prémisses. Une citation authentique est nécessaire à la traçabilité, mais ne suffit pas à prouver une conclusion.

### 39 · [Chapitre 39](../../ch39-active-compliance-auditor-and-popperian-testing.md)

Un auditeur actif recherche des contre-exemples et des cas limites manquants pour les exigences. La génération de tests est séparée de l’autorité qui fixe les résultats attendus ; la couverture de conformité est vérifiable, pas une preuve automatique.

### 40 · [Chapitre 40](../../ch40-distributed-epistemic-architectures-soa-and-cluster-scaling.md)

Le chapitre décrit des systèmes experts distribués avec routage sémantique, ensemble complet de règles actives, coordination de la mémoire et journal des décisions. Un fragment manquant ou une version périmée doit produire un état non résolu, et non une approbation silencieuse.
