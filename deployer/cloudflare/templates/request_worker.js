{{ .Shared }}
{{ .IttyRouter }}

addEventListener('fetch', event => {
  event.respondWith(handleRequest(event.request))
})

async function handleRequest(request) {
  const router = Router()

  try {
    router.{{ .Method }}("{{ .Path }}", async ({ params }) => {
      const body = await request.text();
      return await runRunnable(body, params, request.method.toLowerCase());
    });

    router.all('*', () => new Response('Not Found.', { status: 404 }))

    return router.handle(request);
  } catch (err) {
    return new Response(err.stack, { status: 500 })
  }
}
