import { defineConfig } from 'vitepress'
import { fileURLToPath, URL } from 'node:url'


// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "NutrixPOS | Restaurants & Retail management system",
  description: "The Piece of Cake Restaurant management",
  ignoreDeadLinks: true,
  themeConfig: {
    search: {
      provider: 'local'
    },
  },
  locales: {
    root: {
      label: 'English',
      lang: 'en',
      dir: 'ltr',
      themeConfig: {
        // https://vitepress.dev/reference/default-theme-config
        nav: [
          { text: 'Home', link: '/' },
          { text: 'Downloads', link: 'https://github.com/nutrixpos/pos/releases' },
          { text: 'Docs', link: '/userguide/installation.html' },
          { text: 'About', link: '/about.html' },
          { text: 'Support', link: '/support.html' },
          { text: 'Sponsor', link: '/sponsor.html' },
        ],

        sidebar: [
          {
            text: 'User guide',
            items: [
              { text: 'Installation', link: '/userguide/installation' },
            ]
          },
          {
            text: 'Documentation',
            items: [
              { text: 'Architecture', link: '/architecture' },
              { text: 'docker-compose', link: '/getting_started' },
            ]
          }
        ],

        socialLinks: [
          { icon: 'github', link: 'https://github.com/nutrixpos' }
        ],

        logo: { light: '/logo.png', dark: 'logo-inverse.png', alt: 'logo' },

        search: {
          provider: 'local'
        },

        footer: {
          message: 'Released under the GPL-2.0 License.',
          copyright: 'Copyright © 2026 <a href="https://linkedin.com/in/amr-el-mawardy">Amr Elmawardy</a>'
        }
      },
    },
    ar: {
      label: 'عربي',
      lang: 'ar', // optional, will be added  as `lang` attribute on `html` tag
      link: '/ar',
      dir: 'rtl',
      themeConfig: {
        // https://vitepress.dev/reference/default-theme-config
        nav: [
          { text: 'الرئيسية', link: '/ar' },
          { text: 'تنزيل', link: 'https://github.com/nutrixpos/pos/releases' },
          { text: 'الشرح', link: '/ar/userguide/installation.html' },
          { text: 'عن المشروع', link: '/ar/about.html' },
          { text: 'الدعم الفني', link: '/ar/support.html' },
          { text: 'ادعمنا', link: '/ar/sponsor.html' },
        ],

        sidebar: [
          {
            text: 'دليل المستخدم',
            items: [
              { text: 'التثبيت', link: '/ar/userguide/installation' },
            ]
          },
          {
            text: 'الشرح',
            items: [
              { text: 'التصميم', link: '/ar/architecture' },
              { text: 'docker-compose', link: '/ar/getting_started' },
            ]
          }
        ],

        socialLinks: [
          { icon: 'github', link: 'https://github.com/nutrixpos' }
        ],

        logo: { light: '/logo.png', dark: 'logo-inverse.png', alt: 'logo' },

        search: {
          provider: 'local'
        },

        footer: {
          message: 'مرخص تحت رخص GPL-2.0',
          copyright: ' حقوق النشر  2026 <a href="https://linkedin.com/in/amr-el-mawardy">عمر الماوردي</a>'
        }
      },
    },
    ru: {
      label: 'Русский',
      lang: 'ru', // optional, will be added  as `lang` attribute on `html` tag
      link: '/ru',
      dir: 'ltr',
      themeConfig: {
        // https://vitepress.dev/reference/default-theme-config
        nav: [
          { text: 'Главная', link: '/ru' },
          { text: 'Скачать', link: 'https://github.com/nutrixpos/pos/releases' },
          { text: 'Документация', link: '/ru/userguide/installation.html' },
          { text: 'О нас', link: '/ru/about.html' },
          { text: 'Поддерживать', link: '/ru/support.html' },
          { text: 'Спонсор', link: '/ru/sponsor.html' },
        ],

        sidebar: [
          {
            text: 'Руководство пользователя',
            items: [
              { text: 'Установка', link: '/ru/userguide/installation' },
            ]
          },
          {
            text: 'Документация',
            items: [
              { text: 'Архитектура', link: '/ru/architecture' },
              { text: 'docker-compose', link: '/ru/getting_started' },
            ]
          }
        ],

        socialLinks: [
          { icon: 'github', link: 'https://github.com/nutrixpos' }
        ],

        logo: { light: '/logo.png', dark: 'logo-inverse.png', alt: 'logo' },

        search: {
          provider: 'local'
        },

        footer: {
          message: 'Распространяется под лицензией GPL-2.0.',
          copyright: 'Авторские права © 2026 <a href="https://linkedin.com/in/amr-el-mawardy">Амро Эльмаварди</a>'
        }
      },
    },
  },
  // vite: {
  //   resolve: {
  //     alias: [
  //       {
  //         find: /^.*\/VPNavBarMenuLink\.vue$/,
  //         replacement: fileURLToPath(
  //           new URL("../components/VPNavBarMenuLink.vue", import.meta.url)
  //         )
  //       },
  //       {
  //         find: /^.*\/VPLink\.vue$/,
  //         replacement: fileURLToPath(
  //           new URL("../components/VPNavBarMenuLink.vue", import.meta.url)
  //         )
  //       }
  //     ]
  //   }
  // },
})
