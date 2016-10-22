#!/usr/bin/env python
"""Manage data files"""

import contextlib
import logging
import os
import subprocess
import tempfile
import time

class DataManagerError(Exception):
    pass

@contextlib.contextmanager
def cwd(path):
    curdir= os.getcwd()
    try:
        os.chdir(path)
        yield
    finally:
        os.chdir(curdir)

logger = logging.getLogger(__name__)

def file_age_in_hours(filename):
    """Return a files age in hours or None if it does not exist
    Not the best api, but since this is used for caching it keeps
    the calling function logic simpler
    """
    if not os.path.exists(filename):
        logger.debug("file_age_in_hours: filename=%s does not exist", filename)
        return None
    sr = os.stat(filename)

    age_seconds = time.time() - sr.st_mtime
    age_hours = age_seconds/60/60
    logger.debug("file_age_in_hours: filename=%s hours=%d", filename, age_hours)
    return age_hours

def maybe_update(name, download_func, output_filename, max_age_in_hours=24):
    age = file_age_in_hours(output_filename)
    if age and age < max_age_in_hours:
        logger.debug("maybe_update: func=%s output_filename=%s max_age_in_hours=%d result=cached age=%d",
            name,  output_filename, max_age_in_hours, age)
        return

    return download_func(output_filename)


def download_asnnames(output_filename):
    logger.info("Downloading asn names and writing to %s", output_filename)
    fn = output_filename + ".new"

    if os.path.exists(fn):
        os.unlink(fn)

    cmd = ["pyasn_util_asnames.py", "-o", fn]
    logger.debug("Executing cmd=%s", cmd)
    subprocess.check_call(cmd)
    os.rename(fn, output_filename)

def download_and_convert(output_filename):
    working_dir = tempfile.mkdtemp()

    fn = output_filename + ".new"
    fn_full_path = os.path.join(working_dir, fn)

    logger.debug("download_and_convert working_dir=%s", working_dir)
    with cwd(working_dir):
        cmd = ["pyasn_util_download.py", "--latest"]
        logger.debug("Executing cmd=%s", cmd)
        subprocess.check_call(cmd)

        filenames_here = os.listdir(".")
        if len(filenames_here) != 1:
            raise DataManagerError("More than one file after pyasn_util_download: files=%s", filenames_here)
        rib = filenames_here[0]

        cmd = ["pyasn_util_convert.py", "--single", rib, fn]
        logger.debug("Executing cmd=%s", cmd)
        subprocess.check_call(cmd)
        os.unlink(rib)

    os.rename(fn_full_path, output_filename)
    os.rmdir(working_dir)

def update_asnnames(output_filename, max_age_in_hours=24):
    return maybe_update("names", download_asnnames, output_filename, max_age_in_hours)

def update_asndb(output_filename, max_age_in_hours=24):
    return maybe_update("db", download_and_convert, output_filename, max_age_in_hours)

if __name__ == "__main__":
    logging.basicConfig(level=logging.DEBUG)
    update_asnnames("asnames.json", 24)
    update_asndb("asn.db", 24)
