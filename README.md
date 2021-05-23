# russian-learns-english

I used this repository to learn and remember English words.  
This repository was created because of learning unknown English words and practice my programing skills.  
I used [Yandex dictionary](https://yandex.com/dev/dictionary/) API for getting word's translations.  
To run this API locally, you must create new API key in [Yandex dictionary](https://yandex.com/dev/dictionary/).  
Copy an example file `api.yml.example`.  
You can define port when running this application. Just pass `--port=8080` as a flag.

## /words

**GET** `/words` - Returns all loaded words with their translations.

---
**GET** `/words/random` - Returns random word without translations.

---
**POST** `/words/check-translation` - Check word translation.

---

```
{
    "id": 5687263876523415,
    "translations": ["зима", "зимний"]
}
```

---

**POST** `/words/upload` - Uploads new word list. Use only `text/plain` files. Field name **words**.  
File example:

```
wedding
nod
tidy
```

