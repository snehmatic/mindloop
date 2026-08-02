const CACHE_NAME = 'mindloop-cache-v1';
const ASSETS_TO_CACHE = [
  '/',
  '/static/css/style.css',
  '/static/css/easymde.min.css',
  '/static/js/lucide.min.js',
  '/static/js/purify.min.js',
  '/static/js/easymde.min.js',
  '/static/js/filter_sort.js',
  '/static/js/command_palette.js',
  '/static/images/logo.svg',
  '/static/favicon_io/favicon.svg',
  '/static/favicon_io/android-chrome-192x192.png',
  '/static/favicon_io/android-chrome-512x512.png'
];

self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => {
      return cache.addAll(ASSETS_TO_CACHE);
    })
  );
  self.skipWaiting();
});

self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys().then((cacheNames) => {
      return Promise.all(
        cacheNames.map((name) => {
          if (name !== CACHE_NAME) {
            return caches.delete(name);
          }
        })
      );
    })
  );
  self.clients.claim();
});

self.addEventListener('fetch', (event) => {
  if (event.request.method !== 'GET') return;
  event.respondWith(
    caches.match(event.request).then((response) => {
      return response || fetch(event.request);
    }).catch(() => {
      // Ignore fallback since offline fallback page is not explicitly requested.
    })
  );
});
