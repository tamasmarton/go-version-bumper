# Version Bumper

Egyszerű Go CLI, ami:
- növeli a `package.json` (vagy automatikusan a `packages.json`) verzióját (`major`/`minor`/`patch`),
- gitbe stagingeli és commitolja a fájlt,
- a commit üzenet sablonja konfigurálható a futtatható mellé tett YAML-ból,
- alapértelmezett commit üzenet: csak a verzió (`"%s"`).

## Telepítés

```
go build -o version-bumper .
```

Tedd a `version-bumper` binárist a PATH-odra, vagy futtasd a repo gyökeréből.

## Használat

Alapértelmezésben `patch` bump és automatikus `package.json`/`packages.json` felismerés:

```
./version-bumper
```

Pár példa:
```
./version-bumper minor
./version-bumper major
./version-bumper -file path/to/package.json patch
./version-bumper -dry-run minor
./version-bumper -commit-template "chore(release): %s"
```

Flag-ek:
- `-file` (alap: `package.json`) – útvonal a JSON-hoz; ha ez nem létezik és van `packages.json`, akkor azt használja.
- `-level` (alap: `patch`) – `patch | minor | major`
- `-no-git` – ne fusson `git add` és `git commit`
- `-commit-template` – commit üzenet sablon felülírása flaggel (fmt.Sprintf minta, `%s` a verzió)
- `-dry-run` – csak kiírja, mi történne

## Konfiguráció

A futtatható mellé helyezett `version-bump.yml` (vagy `version-bump.yaml`) fájlból olvassa a commit üzenet sablont.

Példa:
```yaml
commit_message: "chore(release): %s"
```

Ha nincs config, az alap `"%s"`, tehát csak a verziószám kerül az üzenetbe.

Megjegyzés: a sablon a `fmt.Sprintf`-et használja, így egy `%s` helyére kerül a verzió.

## Megjegyzések

- A JSON visszaírása szépen formázva történik (2 szóköz indent). A kulcsok sorrendje és a formázás változhat az első futtatáskor.
- Elő-/build-szuffixed verziókat is elfogad (pl. `1.2.3-beta.1`), bump után a tiszta `MAJOR.MINOR.PATCH` értékre áll.

## Licenc

MIT