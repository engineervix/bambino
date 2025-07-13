/**
 * Custom updater for Go version strings in serve.go
 * This updater finds and updates version strings in the format "version": "x.y.z"
 */

// Ensure CommonJS globals are available
/* global require, module */

const semver = require('semver');

module.exports = {
  readVersion: function (contents) {
    // Find version pattern: "version": "x.y.z"
    const regex = /"version":\s*"([^"]+)"/;
    const match = contents.match(regex);

    if (match && match[1]) {
      const version = match[1];
      // Validate it's a proper semver
      if (semver.valid(version)) {
        return version;
      }
      // If not valid semver, throw a more descriptive error
      throw new Error(`Found version "${version}" but it's not a valid semantic version`);
    }

    throw new Error('Unable to find version pattern "version": "x.y.z" in Go file');
  },

  writeVersion: function (contents, version) {
    // Validate the new version is a proper semver
    if (!semver.valid(version)) {
      throw new Error(`Cannot write invalid version "${version}" - must be valid semantic version`);
    }

    // Replace version pattern: "version": "x.y.z"
    const regex = /"version":\s*"[^"]+"/;
    const newContents = contents.replace(regex, `"version": "${version}"`);

    // Verify the replacement actually happened
    if (newContents === contents) {
      throw new Error('Failed to update version - pattern not found in file');
    }

    return newContents;
  }
};
