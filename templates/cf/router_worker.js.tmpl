{{ .IttyRouter }}

addEventListener('fetch', event => {
  event.respondWith(handleRequest(event.request))
})

async function handleRequest(request) {
  const router = Router()

  try {
  {{- range $key, $value := .WorkersToRoute }}
      router.{{ .Method }}('{{ .Path }}', () => {
          return {{ .ServiceBinding }}.fetch(request);
      });
  {{- end }}

    router.all('*', () => new Response('Not Found.', { status: 404 }))

    return router.handle(request);
  } catch (err) {
    return new Response(err.stack, { status: 500 })
  }
}

