package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gocolly/colly"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type Auto struct {
	Brand string
	Age   string
	Cost  string

	Mileage      string
	Transmission string
	EnginePower  string
	EnginesType  string
	DriveUnit    string

	DromLink string
	IsSell   string
}

var (
	towns = map[string]string{
		"Алтайский край":                    "https://auto.drom.ru/region22/all/",
		"Амурская область":                  "https://auto.drom.ru/region28/all/",
		"Архангельская область":             "https://auto.drom.ru/region29/all/",
		"Астраханская область":              "https://auto.drom.ru/region30/all/",
		"Белгородская область":              "https://auto.drom.ru/region31/all/",
		"Брянская область":                  "https://auto.drom.ru/region32/all/",
		"Владимирская область":              "https://auto.drom.ru/region33/all/",
		"Волгоградская область":             "https://auto.drom.ru/region34/all/",
		"Вологодская область":               "https://auto.drom.ru/region35/all/",
		"Воронежская область":               "https://auto.drom.ru/region36/all/",
		"Еврейская автономная область":      "https://auto.drom.ru/region79/all/",
		"Забайкальский край":                "https://auto.drom.ru/region101/all/",
		"Ивановская область":                "https://auto.drom.ru/region37/all/",
		"Иркутская область":                 "https://auto.drom.ru/region38/all/",
		"Кабардино-Балкарская Республика":   "https://auto.drom.ru/region7/all/",
		"Калининградская область":           "https://auto.drom.ru/region39/all/",
		"Калужская область":                 "https://auto.drom.ru/region40/all/",
		"Камчатский край":                   "https://auto.drom.ru/region41/all/",
		"Карачаево-Черкесская Республика":   "https://auto.drom.ru/region9/all/",
		"Кемеровская область":               "https://auto.drom.ru/region42/all/",
		"Кировская область":                 "https://auto.drom.ru/region43/all/",
		"Костромская область":               "https://auto.drom.ru/region44/all/",
		"Краснодарский край":                "https://auto.drom.ru/region23/all/",
		"Красноярский край":                 "https://auto.drom.ru/region24/all/",
		"Курганская область":                "https://auto.drom.ru/region45/all/",
		"Курская область":                   "https://auto.drom.ru/region46/all/",
		"Ленинградская область":             "https://auto.drom.ru/region47/all/",
		"Липецкая область":                  "https://auto.drom.ru/region48/all/",
		"Магаданская область":               "https://auto.drom.ru/region49/all/",
		"Москва":                            "https://auto.drom.ru/region77/all/",
		"Московская область":                "https://auto.drom.ru/region50/all/",
		"Мурманская область":                "https://auto.drom.ru/region51/all/",
		"Ненецкий автономный округ":         "https://auto.drom.ru/region83/all/",
		"Нижегородская область":             "https://auto.drom.ru/region52/all/",
		"Новгородская область":              "https://auto.drom.ru/region53/all/",
		"Новосибирская область":             "https://auto.drom.ru/region54/all/",
		"Омская область":                    "https://auto.drom.ru/region55/all/",
		"Оренбургская область":              "https://auto.drom.ru/region56/all/",
		"Орловская область":                 "https://auto.drom.ru/region57/all/",
		"Пензенская область":                "https://auto.drom.ru/region58/all/",
		"Пермский край":                     "https://auto.drom.ru/region59/all/",
		"Приморский край":                   "https://auto.drom.ru/region25/all/",
		"Псковская область":                 "https://auto.drom.ru/region60/all/",
		"Республика Адыгея":                 "https://auto.drom.ru/region1/all/",
		"Республика Алтай":                  "https://auto.drom.ru/region4/all/",
		"Республика Башкортостан":           "https://auto.drom.ru/region2/all/",
		"Республика Бурятия":                "https://auto.drom.ru/region3/all/",
		"Республика Дагестан":               "https://auto.drom.ru/region5/all/",
		"Республика Ингушетия":              "https://auto.drom.ru/region6/all/",
		"Республика Калмыкия":               "https://auto.drom.ru/region8/all/",
		"Республика Карелия":                "https://auto.drom.ru/region10/all/",
		"Республика Коми":                   "https://auto.drom.ru/region11/all/",
		"Республика Крым":                   "https://auto.drom.ru/region102/all/",
		"Республика Марий Эл":               "https://auto.drom.ru/region12/all/",
		"Республика Мордовия":               "https://auto.drom.ru/region13/all/",
		"Республика Саха (Якутия)":          "https://auto.drom.ru/region14/all/",
		"Республика Северная Осетия":        "https://auto.drom.ru/region15/all/",
		"Республика Татарстан":              "https://auto.drom.ru/region16/all/",
		"Республика Тыва":                   "https://auto.drom.ru/region17/all/",
		"Республика Хакасия":                "https://auto.drom.ru/region19/all/",
		"Ростовская область":                "https://auto.drom.ru/region61/all/",
		"Рязанская область":                 "https://auto.drom.ru/region62/all/",
		"Самарская область":                 "https://auto.drom.ru/region63/all/",
		"Санкт-Петербург":                   "https://auto.drom.ru/region78/all/",
		"Саратовская область":               "https://auto.drom.ru/region64/all/",
		"Сахалинская область":               "https://auto.drom.ru/region65/all/",
		"Свердловская область":              "https://auto.drom.ru/region66/all/",
		"Смоленская область":                "https://auto.drom.ru/region67/all/",
		"Ставропольский край":               "https://auto.drom.ru/region26/all/",
		"Тамбовская область":                "https://auto.drom.ru/region68/all/",
		"Тверская область":                  "https://auto.drom.ru/region69/all/",
		"Томская область":                   "https://auto.drom.ru/region70/all/",
		"Тульская область":                  "https://auto.drom.ru/region71/all/",
		"Тюменская область":                 "https://auto.drom.ru/region72/all/",
		"Удмуртская Республика":             "https://auto.drom.ru/region18/all/",
		"Ульяновская область":               "https://auto.drom.ru/region73/all/",
		"Хабаровский край":                  "https://auto.drom.ru/region27/all/",
		"Ханты-Мансийский автономный округ": "https://auto.drom.ru/region86/all/",
		"Челябинская область":               "https://auto.drom.ru/region74/all/",
		"Чеченская Республика":              "https://auto.drom.ru/region20/all/",
		"Чувашская Республика":              "https://auto.drom.ru/region21/all/",
		"Чукотский автономный округ":        "https://auto.drom.ru/region87/all/",
		"Ямало-Ненецкий автономный округ":   "https://auto.drom.ru/region89/all/",
		"Ярославская область":               "https://auto.drom.ru/region76/all/",
	}
	currentFile     string = "products.json"
	newFileName     string = "productsNew.json"
	numericKeyboard        = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Следующая машина " + "\u27A1"),
			//tgbotapi.NewKeyboardButton("💔"),
		),
	)
)

func main() {
	bot, err := tgbotapi.NewBotAPI("")
	if err != nil {
		log.Panic(err)
	}

	autos := make([]Auto, 0)
	//auto := Auto{}
	bot.Debug = true
	autos = readerJson(autos, currentFile)
	//fmt.Println(autos[1])
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	i := 0
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		//mysql_db.InsertUserInfo(int(update.Message.Chat.ID))
		var (
			url_town      string = "https://auto.drom.ru/region50/all/"
			selectArea    bool   = false
			toStart       bool   = false
			notifications bool   = false
		)

		if selectArea == false && notifications == false {
			go webScraper(&url_town)
		}

		if update.Message != nil { // If we got a message
			autos = readerJson(autos, currentFile)
			if selectArea == false && notifications == false {
				switch update.Message.Text {
				case "Следующая машина \u27A1":
					i = compareFiles(autos[i], currentFile, newFileName)
					autos = readerJson(autos, currentFile)
					if i != len(autos) {
						if i == 0 && toStart == true {
							i++
							toStart = false
						} else if i != 0 && toStart == false {
							i++
						}
						log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
						replyMain := string("Марка: " + autos[i].Brand + "\n" + "Год: " + autos[i].Age + "\n" + "Пробег: " + autos[i].Mileage + "\n\n")
						replyCosts := string("Цена: " + autos[i].Cost + "\n" + "Статус: " + autos[i].IsSell + "\n" + "Ссылка: " + autos[i].DromLink + "\n\n")
						replyEngine := string("Привод: " + autos[i].EnginesType + "\n" + "КПП: " + autos[i].Transmission + "\n" + "Мощность двигателя: " + autos[i].EnginePower + "\n" + "Тип двигателя: " + autos[i].DriveUnit)
						reply := replyMain + replyCosts + replyEngine
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
						msg.ReplyToMessageID = update.Message.MessageID
						msg.ReplyMarkup = numericKeyboard
						bot.Send(msg)
						if i == 0 && toStart == false {
							i++
						}
					} else if i == len(autos) {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ты просмотрел все машины за сегодня. Жди обновлений. А пока можешь посмотреть список заново :)")
						msg.ReplyToMessageID = update.Message.MessageID
						bot.Send(msg)
						i = 0
					}
				case "/help":
					reply1 := "Для начала просмотра машин нажми на 'Следующая машина \u27A1 \n\n"
					reply2 := "Чтобы вернуться к началу списка напиши команду /to_start \n\n"
					reply3 := "Чтобы задать параметры фильтрации напиши команду /filter (будет реализовано в будущем)\n\n"
					reply4 := "Чтобы обновить список машин и получить актуальную информацию напиши команду /update (будет реализовано в будущем\n\n"
					reply := reply1 + reply2 + reply3 + reply4
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
					msg.ReplyMarkup = numericKeyboard
					bot.Send(msg)
				case "/to_start":
					i = compareFiles(autos[i], currentFile, newFileName)
					autos = readerJson(autos, currentFile)
					i = 0
					toStart = true
					reply := "Возвращаемся к началу списка.\n"
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
					msg.ReplyMarkup = numericKeyboard
					bot.Send(msg)
					if i != len(autos) {
						//i++
						log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
						replyMain := string("Марка: " + autos[i].Brand + "\n" + "Год: " + autos[i].Age + "\n" + "Пробег: " + autos[i].Mileage + "\n\n")
						replyCosts := string("Цена: " + autos[i].Cost + "\n" + "Статус: " + autos[i].IsSell + "\n" + "Ссылка: " + autos[i].DromLink + "\n\n")
						replyEngine := string("Привод: " + autos[i].EnginesType + "\n" + "КПП: " + autos[i].Transmission + "\n" + "Мощность двигателя: " + autos[i].EnginePower + "\n" + "Тип двигателя: " + autos[i].DriveUnit)
						reply = replyMain + replyCosts + replyEngine
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, reply)
						msg.ReplyToMessageID = update.Message.MessageID
						//msg.ReplyMarkup = numericKeyboard
						bot.Send(msg)
					} else if i == len(autos) {
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ты просмотрел все машины за сегодня. Жди обновлений. А пока можешь посмотреть список заново :)")
						msg.ReplyToMessageID = update.Message.MessageID
						bot.Send(msg)
						i = 0
					}
				case "/start":
					reply1 := "Привет, я бот, позволяющий получать информацию о машине. Если хочешь получить моё описание, то введи /help \n\n"
					reply2 := "Для начала просмотра машин нажми на 'Следующая машина \u27A1 \n\n"
					reply := reply1 + reply2
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
					msg.ReplyMarkup = numericKeyboard
					bot.Send(msg)
				case "/select_area":
					selectArea = true
					reply := "Достпуный список регионов: \n\n"
					var town string
					for key, _ := range towns {
						town += key + "\n"
					}
					//msg.ReplyMarkup = numericKeyboard
					reply_end := "\nВыбери интересующий, написав название в чат. Если передумал, то напиши 'отмена' \n"
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply+town+reply_end)
					//
					bot.Send(msg)
				case "/notice_on":
					notifications = true
					reply := "Режим обновлений включен. Чтобы отключить напиши в чат 'notice_off'"
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
					//
					bot.Send(msg)
				default:
					reply := "Я не знаю такой команды"
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
					msg.ReplyMarkup = numericKeyboard
					bot.Send(msg)
				}
			} else if selectArea == true && notifications == false {
				if val, ok := towns[update.Message.Text]; ok {
					url_town = val
					reply := "Город изменён. Обновление может занять некоторое время. Объявление скоро появится..."
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
					//msg.ReplyMarkup = numericKeyboard
					bot.Send(msg)
					webScraper_min(&url_town)
					i = compareFiles(autos[i], currentFile, newFileName)
					autos = readerJson(autos, currentFile)
					if i == 0 && toStart == true {
						i++
						toStart = false
					} else if i != 0 && toStart == false {
						i++
					}
					log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
					replyMain := string("Марка: " + autos[i].Brand + "\n" + "Год: " + autos[i].Age + "\n" + "Пробег: " + autos[i].Mileage + "\n\n")
					replyCosts := string("Цена: " + autos[i].Cost + "\n" + "Статус: " + autos[i].IsSell + "\n" + "Ссылка: " + autos[i].DromLink + "\n\n")
					replyEngine := string("Привод: " + autos[i].EnginesType + "\n" + "КПП: " + autos[i].Transmission + "\n" + "Мощность двигателя: " + autos[i].EnginePower + "\n" + "Тип двигателя: " + autos[i].DriveUnit)
					reply = replyMain + replyCosts + replyEngine
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, reply)
					msg.ReplyToMessageID = update.Message.MessageID
					msg.ReplyMarkup = numericKeyboard
					bot.Send(msg)
					if i == 0 && toStart == false {
						i++
					}
					selectArea = false
				} else if update.Message.Text == "отмена" || update.Message.Text == "Отмена" {
					if selectArea == true {
						selectArea = false
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выбор города отменён")
						bot.Send(msg)
						msg.ReplyMarkup = numericKeyboard
					} else {
						reply := "Я не знаю такой команды"
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
						msg.ReplyMarkup = numericKeyboard
						bot.Send(msg)
					}
				} else {
					reply := "Такой город не найден, попробуй ещё раз или нажми 'отмена'"
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
					bot.Send(msg)
				}
			} else if selectArea == false && notifications == true {
				if update.Message != nil {
					if update.Message.Text == "/notice_off" {
						reply := "Обновления отключены"
						notifications = false
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
						msg.ReplyMarkup = numericKeyboard
						bot.Send(msg)
						break
					} else {
						webScraper_min(&url_town)

						oldAutos := make([]Auto, 0)
						newAutos := make([]Auto, 0)

						oldFile, err := os.Open(currentFile)
						if err != nil {
							log.Fatal(err)
						}

						defer oldFile.Close()

						oldData, err := ioutil.ReadAll(oldFile)

						jsonErrOld := json.Unmarshal(oldData, &oldAutos)

						if err != nil {
							log.Fatal(err)
						}

						if jsonErrOld != nil {
							log.Fatal(jsonErrOld)
						}

						newFile, err := os.Open(newFileName)
						if err != nil {
							log.Fatal(err)
						}

						defer newFile.Close()

						newData, err := ioutil.ReadAll(newFile)

						if err != nil {
							log.Fatal(err)
						}
						jsonErrNew := json.Unmarshal(newData, &newAutos)

						if jsonErrNew != nil {
							log.Fatal(jsonErrNew)
						}
						for j := range newAutos {
							if oldAutos[0].DromLink == newAutos[j].DromLink {
								break
							} else {
								replyMain := string("Марка: " + newAutos[j].Brand + "\n" + "Год: " + newAutos[j].Age + "\n" + "Пробег: " + newAutos[j].Mileage + "\n\n")
								replyCosts := string("Цена: " + newAutos[j].Cost + "\n" + "Статус: " + newAutos[j].IsSell + "\n" + "Ссылка: " + newAutos[j].DromLink + "\n\n")
								replyEngine := string("Привод: " + newAutos[j].EnginesType + "\n" + "КПП: " + newAutos[j].Transmission + "\n" + "Мощность двигателя: " + newAutos[j].EnginePower + "\n" + "Тип двигателя: " + newAutos[j].DriveUnit)
								reply := replyMain + replyCosts + replyEngine
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
								//msg.ReplyMarkup = numericKeyboard
								bot.Send(msg)
							}
						}
						js, err := json.MarshalIndent(newAutos, "", "    ")
						if err != nil {
							log.Fatal(err)
						}
						fmt.Println("Writing data to file")
						if err = os.WriteFile("products.json", js, 0664); err == nil {
							fmt.Println("Data written to file successfully")
						}
					}
				}
			}
		}
	}
}

func readerJson(autos []Auto, file string) []Auto {
	filename, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	defer filename.Close()

	data, err := ioutil.ReadAll(filename)

	if err != nil {
		log.Fatal(err)
	}
	jsonErr := json.Unmarshal(data, &autos)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	//fmt.Println(autos)
	return autos
}

func webScraper(url *string) {
	for {
		// Instantiate default collector
		c := colly.NewCollector(
		//colly.AllowedDomains("coingecko.com"),
		)
		autos := make([]Auto, 0)
		auto := Auto{}
		//flats := make([]Flat, 0)
		c.OnHTML("a.css-5l099z.ewrty961", func(e *colly.HTMLElement) {
			link := e.Attr("href")
			//fmt.Println(link)
			auto.DromLink = link
			e.ForEach("div.css-13ocj84.e727yh30", func(i int, z *colly.HTMLElement) {
				z.ForEach("div", func(i int, u *colly.HTMLElement) {
					if u.ChildText("div.css-r91w5p.e3f4v4l2") == "" {
						u.ForEach("div.css-17lk78h.e3f4v4l2", func(i int, v *colly.HTMLElement) {
							title := v.ChildText("span")
							nameAndDate := strings.Split(title, ", ")
							auto.Brand = nameAndDate[0]
							auto.Age = nameAndDate[1]
							auto.IsSell = "Продаётся"
						})
					} else {
						u.ForEach("div.css-r91w5p.e3f4v4l2", func(i int, v *colly.HTMLElement) {
							title := v.ChildText("span")
							nameAndDate := strings.Split(title, ", ")
							auto.Brand = nameAndDate[0]
							auto.Age = nameAndDate[1]
							auto.IsSell = "Снят с прожажи"
						})
					}
				})
				z.ForEach("div.css-1fe6w6s.e162wx9x0", func(i int, v *colly.HTMLElement) {
					description := v.ChildText("span.css-1l9tp44.e162wx9x0")
					//description := v.ChildText("span.css-1l9tp44.e162wx9x0")
					//description = strings.ReplaceAll(description, " ", "")
					descriptionArr := strings.Split(description, ",")

					for j := range descriptionArr {
						if strings.Contains(descriptionArr[j], "л.с") {
							auto.EnginePower = descriptionArr[j]
						} else if strings.Contains(descriptionArr[j], "АКПП") || strings.Contains(descriptionArr[j], "механика") || strings.Contains(descriptionArr[j], "робот") || strings.Contains(descriptionArr[j], "вариатор") {
							auto.Transmission = descriptionArr[j]
						} else if strings.Contains(descriptionArr[j], "бензин") || strings.Contains(descriptionArr[j], "дизель") || strings.Contains(descriptionArr[j], "гибрид") {
							auto.DriveUnit = descriptionArr[j]
						} else if strings.Contains(descriptionArr[j], "передний") || strings.Contains(descriptionArr[j], "задний") || strings.Contains(descriptionArr[j], "4DW") {
							auto.EnginesType = descriptionArr[j]
						} else if strings.Contains(descriptionArr[j], "тыс. км") {
							descriptionArr[j] = strings.ReplaceAll(descriptionArr[j], "\u003c", "менее")
							auto.Mileage = descriptionArr[j]
						}

					}
				})
			})
			e.ForEach("div.css-1dkhqyq.ep0qbyc0", func(i int, z *colly.HTMLElement) {
				z.ForEach("div", func(i int, u *colly.HTMLElement) {
					u.ForEach("div.css-1i8tk3y.eyvqki92", func(i int, v *colly.HTMLElement) {
						v.ForEach("div.css-1dv8s3l.eyvqki91", func(i int, p *colly.HTMLElement) {
							p.ForEach("span.css-46itwz.e162wx9x0", func(i int, r *colly.HTMLElement) {
								costs := r.Text
								costs = strings.ReplaceAll(costs, "\u00a0", "")
								//fmt.Println(costs)
								auto.Cost = costs
								//fmt.Println(costs)
							})
						})
					})
				})
			})

			autos = append(autos, auto)
		})
		c.OnHTML("a.css-4gbnjj.e24vrp30", func(e *colly.HTMLElement) {
			nextPage := e.Request.AbsoluteURL(e.Attr("href"))
			c.Visit(nextPage)
		})
		c.OnScraped(func(r *colly.Response) {
			fmt.Println("Finished", r.Request.URL)
			js, err := json.MarshalIndent(autos, "", "    ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Writing data to file")
			if err = os.WriteFile(newFileName, js, 0664); err == nil {
				fmt.Println("Data written to file successfully")
			}

		})

		c.OnResponse(func(r *colly.Response) {
			fmt.Println(r.StatusCode)
		})

		// Before making a request print "Visiting ..."
		numVisited := 0
		// Before making a request print "Visiting ..."
		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL.String())
			if numVisited > 30 {
				r.Abort()
			}
			numVisited++
		})

		c.Limit(&colly.LimitRule{
			// Filter domains affected by this rule
			DomainGlob: "*",
			// Set a delay between requests to these domains
			Delay: (1 * time.Second) / 3,
			// Add an additional random delay
			//RandomDelay: 1 * time.Second,
		})

		c.Visit(*url)

		time.Sleep((180 * time.Second) / 2)
		//i = compareFiles(currentFile, newFileName)
	}
}

func compareFiles(auto Auto, currentFile string, newFile string) int {

	newAutos := make([]Auto, 0)
	newFileName, err := os.Open(newFile)
	if err != nil {
		log.Fatal(err)
	}

	defer newFileName.Close()

	newData, err := ioutil.ReadAll(newFileName)

	if err != nil {
		log.Fatal(err)
	}
	jsonErrNew := json.Unmarshal(newData, &newAutos)

	if jsonErrNew != nil {
		log.Fatal(jsonErrNew)
	}
	var k int = 0
	for j := range newAutos {
		if auto.DromLink == newAutos[j].DromLink {
			k = j
			break
		}
	}
	js, err := json.MarshalIndent(newAutos, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Writing data to file")
	if err = os.WriteFile("products.json", js, 0664); err == nil {
		fmt.Println("Data written to file successfully")
	}
	return k
}

func webScraper_min(url *string) {
	// Instantiate default collector
	d := colly.NewCollector(
	//colly.AllowedDomains("coingecko.com"),
	)
	autos := make([]Auto, 0)
	auto := Auto{}
	//flats := make([]Flat, 0)
	d.OnHTML("a.css-5l099z.ewrty961", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		//fmt.Println(link)
		auto.DromLink = link
		e.ForEach("div.css-13ocj84.e727yh30", func(i int, z *colly.HTMLElement) {
			z.ForEach("div", func(i int, u *colly.HTMLElement) {
				if u.ChildText("div.css-r91w5p.e3f4v4l2") == "" {
					u.ForEach("div.css-17lk78h.e3f4v4l2", func(i int, v *colly.HTMLElement) {
						title := v.ChildText("span")
						nameAndDate := strings.Split(title, ", ")
						auto.Brand = nameAndDate[0]
						auto.Age = nameAndDate[1]
						auto.IsSell = "Продаётся"
					})
				} else {
					u.ForEach("div.css-r91w5p.e3f4v4l2", func(i int, v *colly.HTMLElement) {
						title := v.ChildText("span")
						nameAndDate := strings.Split(title, ", ")
						auto.Brand = nameAndDate[0]
						auto.Age = nameAndDate[1]
						auto.IsSell = "Снят с прожажи"
					})
				}
			})
			z.ForEach("div.css-1fe6w6s.e162wx9x0", func(i int, v *colly.HTMLElement) {
				description := v.ChildText("span.css-1l9tp44.e162wx9x0")
				//description := v.ChildText("span.css-1l9tp44.e162wx9x0")
				//description = strings.ReplaceAll(description, " ", "")
				descriptionArr := strings.Split(description, ",")

				for j := range descriptionArr {
					if strings.Contains(descriptionArr[j], "л.с") {
						auto.EnginePower = descriptionArr[j]
					} else if strings.Contains(descriptionArr[j], "АКПП") || strings.Contains(descriptionArr[j], "механика") || strings.Contains(descriptionArr[j], "робот") || strings.Contains(descriptionArr[j], "вариатор") {
						auto.Transmission = descriptionArr[j]
					} else if strings.Contains(descriptionArr[j], "бензин") || strings.Contains(descriptionArr[j], "дизель") || strings.Contains(descriptionArr[j], "гибрид") {
						auto.DriveUnit = descriptionArr[j]
					} else if strings.Contains(descriptionArr[j], "передний") || strings.Contains(descriptionArr[j], "задний") || strings.Contains(descriptionArr[j], "4DW") {
						auto.EnginesType = descriptionArr[j]
					} else if strings.Contains(descriptionArr[j], "тыс. км") {
						descriptionArr[j] = strings.ReplaceAll(descriptionArr[j], "\u003c", "менее")
						auto.Mileage = descriptionArr[j]
					}

				}
			})
		})
		e.ForEach("div.css-1dkhqyq.ep0qbyc0", func(i int, z *colly.HTMLElement) {
			z.ForEach("div", func(i int, u *colly.HTMLElement) {
				u.ForEach("div.css-1i8tk3y.eyvqki92", func(i int, v *colly.HTMLElement) {
					v.ForEach("div.css-1dv8s3l.eyvqki91", func(i int, p *colly.HTMLElement) {
						p.ForEach("span.css-46itwz.e162wx9x0", func(i int, r *colly.HTMLElement) {
							costs := r.Text
							costs = strings.ReplaceAll(costs, "\u00a0", "")
							//fmt.Println(costs)
							auto.Cost = costs
							//fmt.Println(costs)
						})
					})
				})
			})
		})

		autos = append(autos, auto)
	})
	d.OnHTML("a.css-4gbnjj.e24vrp30", func(e *colly.HTMLElement) {
		nextPage := e.Request.AbsoluteURL(e.Attr("href"))
		d.Visit(nextPage)
	})
	d.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
		js, err := json.MarshalIndent(autos, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Writing data to file")
		if err = os.WriteFile(newFileName, js, 0664); err == nil {
			fmt.Println("Data written to file successfully")
		}

	})

	d.OnResponse(func(r *colly.Response) {
		fmt.Println(r.StatusCode)
	})

	// Before making a request print "Visiting ..."
	numVisited := 0
	// Before making a request print "Visiting ..."
	d.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
		if numVisited > 5 {
			r.Abort()
		}
		numVisited++
	})

	d.Limit(&colly.LimitRule{
		// Filter domains affected by this rule
		DomainGlob: "*",
		// Set a delay between requests to these domains
		Delay: (1 * time.Second) / 2,
		// Add an additional random delay
		//RandomDelay: 1 * time.Second,
	})

	d.Visit(*url)

	//time.Sleep((180 * time.Second) / 2)
	//i = compareFiles(currentFile, newFileName)
}
