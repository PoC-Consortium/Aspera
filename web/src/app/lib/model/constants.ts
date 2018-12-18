/*
* Copyright 2018 PoC-Consortium
*/

export const constants = {
    connectionTimeout: 10000,
    currencies: [
        { code:  "AUD", symbol: "$" },
        { code:  "BRL", symbol: "R$" },
        { code:  "CAD", symbol: "$" },
        { code:  "CLP", symbol: "$" },
        { code:  "CNY", symbol: "¥" },
        { code:  "CZK", symbol: "Kč" },
        { code:  "DKK", symbol: "kr" },
        { code:  "EUR", symbol: "€" },
        { code:  "GBP", symbol: "£" },
        { code:  "HKD", symbol: "$" },
        { code:  "HUF", symbol: "Ft" },
        { code:  "IDR", symbol: "Rp" },
        { code:  "ILS", symbol: "₪" },
        { code:  "INR", symbol: "₹" },
        { code:  "JPY", symbol: "¥" },
        { code:  "KRW", symbol: "₩" },
        { code:  "MXN", symbol: "$" },
        { code:  "MYR", symbol: "RM" },
        { code:  "NOK", symbol: "kr" },
        { code:  "NZD", symbol: "$" },
        { code:  "PHP", symbol: "₱" },
        { code:  "PKR", symbol: "₨" },
        { code:  "PLN", symbol: "zł" },
        { code:  "RUB", symbol: "₽" },
        { code:  "SEK", symbol: "kr" },
        { code:  "SGD", symbol: "$" },
        { code:  "THB", symbol: "฿" },
        { code:  "TRY", symbol: "₺" },
        { code:  "TWD", symbol: "$" },
        { code:  "USD", symbol: "$" },
        { code:  "ZAR", symbol: "Rs" }
    ],
    database: "loki.db",
    defaultCurrency: "USD",
    defaultLanguage: "en",
    defaultTheme: "light",
    documentationUrl: "https://poc-consortium.github.io/burstcoin-mobile-doc/",
    donate: "BURST-RTEY-HUSA-BJG4-EZW9E",
    languages: [
        { code: "en", name: "English (Default)" },
        { code: "bg", name: "Български" },
        { code: "ca", name: "Català" },
        { code: "cs", name: "Čeština" },
        { code: "de-de", name:"Deutsch" },
        { code: "el", name:"Ελληνικά" },
        { code: "es-es", name:"Español" },
        { code: "fi", name:"Suomi" },
        { code: "fr", name:"Français" },
        { code: "gl", name:"Galego" },
        { code: "hi", name:"हिंदी" },
        { code: "hr", name:"Hrvatski" },
        { code: "id", name:"Bahasa Indonesia" },
        { code: "it", name:"Italiano" },
        { code: "ja", name:"日本語" },
        { code: "lt", name:"Lietuviškai" },
        { code: "nl", name:"Nederlands" },
        { code: "sh", name:"Hrvatski" },
        { code: "sk", name:"Slovensky" },
        { code: "pt-pt", name:"Português (Portugal)" },
        { code: "pt-br", name:"Português (Brazil)" },
        { code: "sr", name:"Српски" },
        { code: "sr-cs", name:"Srpski" },
        { code: "tr", name:"Türk" },
        { code: "uk", name:"Yкраiнска" },
        { code: "ro", name:"Român" },
        { code: "ru", name:"Русский" },
        { code: "zh-cn", name:"中文 (simplified)" },
        { code: "zh-tw", name:"中文 (traditional)" }
    ],
    nodes: [
        { "name": "CryptoGuru", "region": "Global", "location": "Munich", "address": "https://wallet.burst.cryptoguru.org", "port": 8125, "selected": true, "ping": -1 },
        { "name": "CryptoGuru", "region": "Africa", "location": "Munich", "address": "https://wallet.burst.cryptoguru.org", "port": 8125, "selected": false, "ping": -1 },
        { "name": "CryptoGuru", "region": "Asia", "location": "Munich", "address": "https://wallet.burst.cryptoguru.org", "port": 8125, "selected": false, "ping": -1 },
        { "name": "CryptoGuru", "region": "Europe", "location": "Munich", "address": "https://wallet.burst.cryptoguru.org", "port": 8125, "selected": false, "ping": -1 },
        { "name": "CryptoGuru", "region": "North America", "location": "Munich", "address": "https://wallet.burst.cryptoguru.org", "port": 8125, "selected": false, "ping": -1 },
        { "name": "CryptoGuru", "region": "Oceania", "location": "Munich", "address": "https://wallet.burst.cryptoguru.org", "port": 8125, "selected": false, "ping": -1 },
        { "name": "CryptoGuru", "region": "South America", "location": "Munich", "address": "https://wallet.burst.cryptoguru.org", "port": 8125, "selected": false, "ping": -1 },
    ],
    supportUrl: "https://github.com/poc-consortium/burstcoin-mobile-doc/issues",
    transactionCount: "500",
    transactionUrl: "https://explore.burst.cryptoguru.org/transaction/",
    twitter: "https://twitter.com/PoC_Consortium"
}
