import lines from "./in/pages.json";
import segments from "./in/qpc-v4.json";
import uthmani from "./in/uthmani.json";
import fs from "fs/promises";

type segment = {
  id: number;
  surah: string;
  ayah: string;
  word: string;
  location: `${number}:${number}:${number}`;
  text: string;
};

type ayah = {
  surah: number;
  ayah: number;
  glyphs: string;
  text: string;
  page: number;
};

const surahs: { [key: number]: ayah[] } = {};

const ayahs: {
  [key: string]: ayah;
} = {};

const words: {
  [key: number]: {
    surah: number;
    ayah: number;
    word: number;
    location: string;
    text: string;
    page: number;
  };
} = {};

Object.keys(segments).forEach((key) => {
  const word = (segments as { [key: string]: segment })[key];

  if (word) {
    words[word.id] = {
      surah: Number(word.surah),
      word: Number(word.word),
      ayah: Number(word.ayah),
      location: word.location,
      text: word.text,
      page: -1,
    };
  }
});

lines.forEach((line) => {
  if (line.line_type == "ayah") {
    for (
      let index = line.first_word_id as number;
      index <= (line.last_word_id as number);
      index++
    ) {
      const word = words[index];
      if (word != undefined) {
        word.page = line.page_number;
      }
    }
  }
});

Object.keys(segments).forEach((key) => {
  const word = (segments as { [key: string]: segment })[key];
  const uthmaniWord = (uthmani as { [key: string]: segment })[key];
  if (word != undefined && uthmaniWord != undefined) {
    const [surahIndex, ayahIndex, wordIndex] = key.split(":");
    const ayahKey = `${surahIndex}:${ayahIndex}`;

    if (ayahs[ayahKey] != undefined) {
      ayahs[ayahKey]["glyphs"] += ` ${word.text}`;
      ayahs[ayahKey]["text"] += ` ${uthmaniWord.text}`;
    } else {
      ayahs[ayahKey] = {
        ayah: Number(ayahIndex),
        surah: Number(surahIndex),
        glyphs: word.text,
        text: uthmaniWord.text,
        page: (words[word.id] as (typeof words)[number]).page,
      };
    }
  }
});

Object.keys(ayahs).forEach((key) => {
  const ayah = ayahs[key];
  if (ayah) {
    if (surahs[ayah.surah] == undefined) {
      surahs[ayah.surah] = [];
    }
    (surahs[ayah.surah] as ayah[]).push(ayah);
  }
});

Object.keys(surahs).forEach((key) => {
  const surah = surahs[Number(key)];

  if (surah) {
    surah.sort((a, b) => a.ayah - b.ayah);
  }
});

await fs.writeFile("./out/ayahs.json", JSON.stringify(ayahs, null, 2));
await fs.writeFile("./out/surahs.json", JSON.stringify(surahs, null, 2));
