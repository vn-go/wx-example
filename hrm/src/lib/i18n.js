import { addMessages, init, locale } from 'svelte-i18n';

export async function loadLocale(lang) {
    const res = await fetch(`/api/translations?lang=${lang}`);
    if (!res.ok) {
        throw new Error(`Failed to load translations for ${lang}`);
    }

    const messages = await res.json();
    addMessages(lang, messages);
    locale.set(lang);
}

init({
    fallbackLocale: 'en',
    initialLocale: 'en'
});
