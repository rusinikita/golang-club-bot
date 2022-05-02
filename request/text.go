package request

import (
	"github.com/rusinikita/gogoClub/simplectx"
)

// requests

const (
	hello    simplectx.Sticker = "CAACAgQAAxkBAAOpYmmXKS3ykycZk2qrR97R2_jTLKwAAswAA845CA3fZ3xlfkS5ZCQE"   // 🙃
	wtf      simplectx.Sticker = "CAACAgQAAxkBAAOrYmmXgeYtrZiN2IuUckR854EheykAApkAA845CA0jIAABUzXpH78kBA" // 🤨
	wellDone simplectx.Sticker = "CAACAgQAAxkBAAOtYmmX5-cveqGjl44BirOjkuy1cz4AApcAA845CA1AIS58gGBWGiQE"   // 👍
)

const (
	hiText = `Привет!

Если ты хочешь поучаствовать <a href="%s">в клубе изучения Go</a>, то ты по адресу.

Пройди простое задание, чтобы доказать мотивацию, и сможешь учиться GoLang с теми, кто поможет справиться с трудностями и подскажет что делать после.`
	step1 = `1. Создай профиль на github.com

Заполни настоящее имя, фото и город. Напиши в bio текущее место учебы/работы (направление и курс тоже)`
	step2 = `2. <a href="https://drive.google.com/file/d/1-8AQtU5WuftQrUioXYkp0bY2K20t9vM3/view?usp=sharing">Создай репозиторий</a>

Название kgl-go-learing (с дополнениями если занято)
README файл (нужно поставить галочку)

Этот репозиторий станет твоим портфолио, тетрадкой с заданиями и поможет систематизировать знания`
	step3 = `3. Напиши в README файле

1. О себе
2. Цель изучения. Почему ты хочешь научиться go или программированию ("почему бы и нет" - валидный ответ, но нужно подробнее расписать)
3. Почему уверен, что не бросишь занятия и не потратишь время менторов впустую
4. Ожидания. Как скоро и что ты хочешь получить от участия.
5. Перечисли вопросы, которые хотел бы обсудить.

Используй # для создания заголовков`
	step4 = `4. Просто отправь мне ссылку на получившийся репозиторий в сообщении.

Например, https://github.com/rusinikita/mindful-bot`
	step5      = `5. Дождись конца проверки и ответа`
	step6      = `Это всё, жду ссылку. Если какие-то проблемы, предложения, что-то не нравится - напиши <a href="%s">в этом чате</a>`
	noLinkText = `Ты не прислал ссылку. Если хочешь обсудить что-то - напиши <a href="%s">в этом чате</a>`
	doneText   = `Супер. Как только проверят задание, я отправлю тебе результат.

Если хочешь поменять ссылку, отправь новую.`
)

// notifications

const (
	sad     simplectx.Sticker = "CAACAgQAAx0CWX9zlwADBGJvzU00-F8LdLK0nQhgnt0JQqtoAAKoAAPOOQgN2hWbG1Xxf5YkBA"
	welcome simplectx.Sticker = "CAACAgQAAx0CWX9zlwADCGJvzvXOyCeUSrocMfNHBjLZ1BOeAAJyAAPOOQgNni2XMTfi1_okBA"
	sorry   simplectx.Sticker = "CAACAgQAAx0CWX9zlwADBmJvzdKFpMv9Vj819LJboKkbh9e2AAJFAAPOOQgNHZQ4AAFROQQ9JAQ"
)

const (
	remind = `Привет! Я заинтересовал тебя ранее, но все ещё жду твою заявку.

Уже %d человек приняты в клуб, но у тебя ещё есть шанс.
Всего нас будет 15-20 человек. Новички и опытные программисты.

Будет весело и полезно!`
	decline = `Привет! Мне жаль, но твоя заявка отклонена.`
	accept  = `Привет! Рад буду видеть тебя в нашей группе.

Только давай договоримся о правилах.
Исключение за дискриминацию, насилие или её поддержку, а так же за бессмысленные или неаргументированные споры на любую тему`
)
