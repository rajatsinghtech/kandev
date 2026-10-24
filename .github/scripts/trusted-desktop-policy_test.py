"""Credential-free tests for stable desktop publication admission."""
import importlib.util
import json
from pathlib import Path
import unittest

ROOT = Path(__file__).resolve().parents[2]
SPEC = importlib.util.spec_from_file_location('desktop_policy', ROOT / 'scripts/release/trusted-desktop-policy.py')
POLICY = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(POLICY)


class TrustedDesktopPolicyTests(unittest.TestCase):
    def setUp(self):
        self.env = dict.fromkeys([
            'APPLE_CERTIFICATE', 'APPLE_CERTIFICATE_PASSWORD', 'KEYCHAIN_PASSWORD',
            'APPLE_ID', 'APPLE_PASSWORD', 'APPLE_TEAM_ID', 'WINDOWS_CERTIFICATE',
            'WINDOWS_CERTIFICATE_PASSWORD', 'TAURI_SIGNING_PRIVATE_KEY'], 'present')
        self.repository = 'rajatsinghtech/kandev'
        self.config = {'plugins': {'updater': {
            'pubkey': 'fork-public-key',
            'endpoints': [f'https://github.com/{self.repository}/releases/latest/download/latest.json']}}}

    def test_complete_inputs_admitted(self):
        self.assertEqual(POLICY.publication_errors(self.env, self.repository, self.config), [])

    def test_missing_or_whitespace_input_refused(self):
        for key in self.env:
            for value in ('', '  '):
                env = dict(self.env, **{key: value})
                with self.subTest(key=key, value=value):
                    self.assertTrue(POLICY.publication_errors(env, self.repository, self.config))

    def test_upstream_endpoint_and_key_refused_for_fork(self):
        for field, value in [('endpoints', ['https://github.com/kdlbs/kandev/releases/latest/download/latest.json']),
                             ('pubkey', POLICY.UPSTREAM_PUBLIC_KEY), ('pubkey', '')]:
            config = json.loads(json.dumps(self.config))
            config['plugins']['updater'][field] = value
            self.assertTrue(POLICY.publication_errors(self.env, self.repository, config))

    def test_errors_never_disclose_secret_values(self):
        env = dict(self.env, WINDOWS_CERTIFICATE='')
        self.assertNotIn('present', str(POLICY.publication_errors(env, self.repository, self.config)))

    def test_publication_workflow_enforces_policy_and_artifact_verification(self):
        workflow = (ROOT / '.github/workflows/release.yml').read_text()
        self.assertGreaterEqual(workflow.count('python3 scripts/release/trusted-desktop-policy.py'), 2)
        self.assertIn('Verify macOS publisher identity and notarization', workflow)
        self.assertIn('xcrun stapler validate', workflow)
        self.assertIn('cargo audit --file apps/desktop/src-tauri/Cargo.lock', workflow)
        self.assertNotIn('--allow-unsigned', workflow)


if __name__ == '__main__':
    unittest.main()
