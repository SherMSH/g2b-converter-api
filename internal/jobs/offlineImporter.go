package jobs

import (
	configpkg "converterapi/internal/config"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"strings"
)

// Партнёр выкладывает пакеты на выпуск карт по SFTP, в тот же каталог, откуда
// мы забираем транзакционные файлы. Сканер офлайн-пакетов туда не ходит - он
// разбирает локальный каталог, - поэтому файлы нужно сначала перенести.

// offlinePrefixes - канонические префиксы имён пакетов: CreateCardsOut,
// ReissueCardsOut и остальные из utils.OfflineReqTypes без маски "*.xml".
func offlinePrefixes() []string {
	prefixes := make([]string, 0, len(utils.OfflineReqTypes))
	for _, t := range utils.OfflineReqTypes {
		prefixes = append(prefixes, strings.TrimSuffix(string(t), "*.xml"))
	}
	return prefixes
}

// normalizeName убирает разделители и регистр.
//
// Партнёр называет файлы по-своему: "CreateCards_Out_2026.09.29_16.36.44"
// вместо "CreateCardsOut_...xml". Сравнивать имена буквально значит терять
// такие пакеты молча, поэтому сравниваем по «скелету» имени.
func normalizeName(name string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(name) {
		switch r {
		case '_', '-', '.', ' ':
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// offlineLocalName подбирает канонное имя для пакета: <Префикс><остаток>.xml.
//
// Возвращает пустую строку, если файл не похож ни на один тип пакета.
func offlineLocalName(remote string) string {
	norm := normalizeName(remote)

	for _, prefix := range offlinePrefixes() {
		normPrefix := normalizeName(prefix)
		if !strings.HasPrefix(norm, normPrefix) {
			continue
		}

		// Отрезаем от исходного имени ровно те символы, что съел префикс:
		// разделители внутри него считаем частью префикса
		var consumed, matched int
		for _, r := range remote {
			if matched == len(normPrefix) {
				break
			}
			consumed++
			if normalizeName(string(r)) != "" {
				matched++
			}
		}

		rest := strings.TrimLeft(remote[consumed:], "_-. ")
		rest = strings.TrimSuffix(rest, ".xml")

		name := prefix
		if rest != "" {
			name += "_" + rest
		}
		return name + ".xml"
	}
	return ""
}

// ImportOfflinePackets переносит пакеты выпуска карт с SFTP во входной каталог
// сканера. Разбирает их уже ConvScanner на своём цикле.
func ImportOfflinePackets(remoteDir string) error {
	localDir := configpkg.Config.App.Storage.Basepath + configpkg.Config.App.Storage.In

	// Обработанные пакеты сканер уносит в out и errors - без них файл считался
	// бы новым и скачивался заново на каждом заходе
	tracker, err := NewImportedFilesTracker(
		localDir,
		configpkg.Config.App.Storage.Basepath+configpkg.Config.App.Storage.Out,
		configpkg.Config.App.Storage.Basepath+configpkg.Config.App.Storage.Errors,
	)
	if err != nil {
		return err
	}

	client, err := trnSFTP.client()
	if err != nil {
		return err
	}

	files, err := client.ReadDir(remoteDir)
	if err != nil {
		trnSFTP.drop()
		return err
	}

	var imported int
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		localName := offlineLocalName(file.Name())
		if localName == "" || tracker.isImported(localName) {
			continue
		}

		if err := importSingleFile(client, remoteDir+"/"+file.Name(), localDir, file.Size(), localName); err != nil {
			logger.Errorf("[JOBS] ошибка импорта пакета %s: %v", file.Name(), err)
			continue
		}

		tracker.markAsImported(localName)
		imported++
		logger.Infof("[JOBS] пакет %s импортирован как %s (%d байт)", file.Name(), localName, file.Size())
	}

	// Не ждём своего расписания: пакет уже лежит во входном каталоге
	if imported > 0 {
		ConvScanner()
	}
	return nil
}
