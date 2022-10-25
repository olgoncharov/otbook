#!/usr/bin/env python3

import argparse
from datetime import datetime

import jwt
import requests

from base import Ammo, AmmoFileGenerator


class AmmoProvider:
    def __init__(self, profiles_url, friends_url_pattern, jwt_secret):
        self._profiles_url = profiles_url
        self._friends_url_pattern = friends_url_pattern
        self._jwt_secret = jwt_secret

    def _make_ammo(self, current_user, new_friend):
        access_token = jwt.encode(
            {
                'iat': datetime.now().timestamp(),
                'exp': datetime(year=2999, month=1, day=1).timestamp(),
                'username': current_user['username']
            },
            self._jwt_secret,
            algorithm='HS256'
        )

        headers = (
            'Host: 127.0.0.1:8000\r\n'
            'User-Agent: tank\r\n'
            'Accept: */*\r\n'
            'Connection: Close\r\n'
            f'Authorization: Bearer {access_token}'
        )

        return Ammo(
            method='POST',
            url=self._friends_url_pattern.format(new_friend['username']),
            headers=headers,
            tag='friends_create',
            body=None
        )

    def get(self, ammo_count):
        response = requests.get(self._profiles_url, params={'limit': ammo_count+1}).json()

        return [
            self._make_ammo(response['list'][i], response['list'][i+1])
            for i in range(ammo_count)
        ]


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('-o', '--output', help='Output filename with ammo', required=True, dest='output_filename')
    parser.add_argument('-a', '--amount', help='Required ammo amount', default=1000, dest='ammo_count', type=int)
    args = parser.parse_args()

    provider = AmmoProvider(
        'http://127.0.0.1:8000/api/v1/profiles',
        '/api/v1/profiles/{}/friends',
        'nnUdOVjyvgtsTUPguUrm'
    )

    AmmoFileGenerator(provider, args.output_filename, args.ammo_count).run()
