from dataclasses import dataclass
from os import stat


@dataclass
class Ammo:
    method: str
    url: str
    headers: str
    tag: str
    body: str


class AmmoFileGenerator:
    def __init__(self, provider, filename, ammo_count):
        self._filename = filename
        self._provider = provider
        self._ammo_count = ammo_count

    @staticmethod
    def _make_request(method, url, headers, body=None):
        if not body:
            return (
                f'{method} {url} HTTP/1.1\r\n'
                f'{headers}\r\n'
                '\r\n'
            )

        return (
            f'{method} {url} HTTP/1.1\r\n'
            f'{headers}\r\n'
            f'Content-Length: {len(body)}\r\n'
            '\r\n'
            f'{body}\r\n'
        )

    def run(self):
        with open(self._filename, 'w') as f:
            for ammo in self._provider.get(self._ammo_count):
                req = self._make_request(ammo.method, ammo.url, ammo.headers, ammo.body)
                f.write(f'{len(req)} {ammo.tag}\n')
                f.write(req)
