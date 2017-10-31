import * as CryptoJS from "crypto-js";
import * as BN from "bn.js";

declare function escape(s: string): string;
declare function unescape(s: string): string;

export class Converter {

    // Convert a hex string to a byte array
    public static convertHexStringToByteArray(hex): number[] {
        for (var bytes = [], c = 0; c < hex.length; c += 2)
            bytes.push(parseInt(hex.substr(c, 2), 16));
        return bytes;
    }

    // Convert a byte array to a hex string
    public static convertByteArrayToHexString(bytes) {
        for (var hex = [], i = 0; i < bytes.length; i++) {
            hex.push((bytes[i] >>> 4).toString(16));
            hex.push((bytes[i] & 0xF).toString(16));
        }
        return hex.join("");
    }

    public static convertByteArrayToWordArray(ba) {
        var wa = [],
            i;
        for (i = 0; i < ba.length; i++) {
            wa[(i / 4) | 0] |= ba[i] << (24 - 8 * i);
        }

        return CryptoJS.lib.WordArray.create(wa)
    }

    public static convertWordToByteArray(word, length) {
        var ba = [],
            i,
            xFF = 0xFF;
        if (length > 0)
            ba.push(word >>> 24);
        if (length > 1)
            ba.push((word >>> 16) & xFF);
        if (length > 2)
            ba.push((word >>> 8) & xFF);
        if (length > 3)
            ba.push(word & xFF);

        return ba;
    }

    public static convertWordArrayToByteArray(wordArray, length = 0) {
        if (wordArray.hasOwnProperty("sigBytes") && wordArray.hasOwnProperty("words")) {
            length = wordArray.sigBytes;
            wordArray = wordArray.words;
        }

        let result: number[] = [],
            bytes,
            i = 0;
        while (length > 0) {
            bytes = Converter.convertWordToByteArray(wordArray[i], Math.min(4, length));
            length -= bytes.length;
            result.push(bytes);
            i++;
        }
        return [].concat.apply([], result);
    }

    public static convertTimestampToDateString(timestamp: number) {
        return new Date(Date.UTC(2014, 7, 11, 2, 0, 0, 0) + timestamp * 1000);
    }
}
