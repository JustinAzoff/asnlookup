#!/usr/bin/env python
"""Manage data files"""

import os
import subprocess
import time
import logging

logger = logging.getLogger(__name__)

def download_asnnames(output_filename):
    logger.info("Downloading asn names and writing to %s", output_filename)
    fn = output_filename + ".new"

    if os.path.exists(fn):
        os.unlink(fn)

    cmd = ["pyasn_util_asnames.py", "-o", fn]
    logger.debug("Executing %s", cmd)
    subprocess.check_call(cmd)
    os.rename(fn, output_filename)

def file_age_in_hours(filename):
    """Return a files age in hours or None if it does not exist
    Not the best api, but since this is used for caching it keeps
    the calling function logic simpler
    """
    if not os.path.exists(filename):
        return None
    sr = os.stat(filename)

    age_seconds = time.time() - sr.st_mtime
    return age_seconds/(60*60)

def update_asnnames(output_filename, max_age_in_hours=24):
    age = file_age_in_hours(output_filename)
    if age and age < max_age_in_hours:
        logger.debug("age of %s is %d hours and is less than %d hours, not downloading", output_filename, age, max_age_in_hours)
        return

    download_asnnames(output_filename)

if __name__ == "__main__":
    logging.basicConfig(level=logging.DEBUG)
    update_asnnames("asnames.json", 24)
