"""Admit trusted desktop publication only with complete owner signing inputs."""
import argparse
import json
import os
from pathlib import Path
import re
import sys

UPSTREAM_PUBLIC_KEY = 'dW50cnVzdGVkIGNvbW1lbnQ6IG1pbmlzaWduIHB1YmxpYyBrZXk6IEI0NzVGQ0MwODhDODM0OEMKUldTTU5NaUl3UHgxdENoYmIvRDR1Y3F1Tys2ZXN5MFNRbHRHSHpmeWhIdzIxVFJEVUJubDAwWmwK'


def publication_errors(env, repository, config):
    errors = []
    required = ['APPLE_CERTIFICATE', 'APPLE_CERTIFICATE_PASSWORD', 'KEYCHAIN_PASSWORD',
                'WINDOWS_CERTIFICATE', 'WINDOWS_CERTIFICATE_PASSWORD', 'TAURI_SIGNING_PRIVATE_KEY']
    for key in required:
        if not env.get(key, '').strip():
            errors.append(f'Missing signing input: {key}')
    apple_id = all(env.get(k, '').strip() for k in ['APPLE_ID', 'APPLE_PASSWORD', 'APPLE_TEAM_ID'])
    apple_api = all(env.get(k, '').strip() for k in ['APPLE_API_KEY', 'APPLE_API_ISSUER', 'APPLE_API_KEY_P8'])
    if not apple_id and not apple_api:
        errors.append('Complete Apple notarization credentials are required')
    if not re.fullmatch(r'[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+', repository):
        errors.append('A valid release repository is required')
    updater = config.get('plugins', {}).get('updater', {})
    expected = f'https://github.com/{repository}/releases/latest/download/latest.json'
    if updater.get('endpoints') != [expected]:
        errors.append('Updater endpoint must belong to the release repository')
    public_key = updater.get('pubkey', '').strip()
    if not public_key or repository != 'kdlbs/kandev' and public_key == UPSTREAM_PUBLIC_KEY:
        errors.append('A fork-owned updater public key is required')
    return errors


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--repository', required=True)
    parser.add_argument('--config', default='apps/desktop/src-tauri/tauri.conf.json', type=Path)
    args = parser.parse_args()
    errors = publication_errors(os.environ, args.repository, json.loads(args.config.read_text(encoding='utf-8')))
    if errors:
        print('\n'.join(errors), file=sys.stderr)
        return 1
    print('Desktop signing prerequisites verified; artifact signature verification remains required.')
    return 0


if __name__ == '__main__':
    raise SystemExit(main())
