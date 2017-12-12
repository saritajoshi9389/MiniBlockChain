#  Citation : Online tutorial for FUSE filesystem on Youtube and internet
#!/usr/bin/python
import sys
import stat
import errno
import time
import re
from fuse import FUSE, FuseOSError, Operations
from Custom import Custom


class DirEntry(object):
    def __init__(self):
        pass

    def stat(self):
        now = time.time()
        # print "stat", dict(st_mode=(stat.S_IFDIR | 0o755), st_ctime=now, st_mtime=now, st_atime=now, st_nlink=2)
        return dict(st_mode=(stat.S_IFDIR | 0o755), st_ctime=now, st_mtime=now, st_atime=now, st_nlink=2)


class FileEntry(object):
    def __init__(self, data):
        self.data = data
        if data:self.size = len(data)
        else:
            print "No more data available"
            exit(2)

    def stat(self):
        now = time.time()
        return dict(st_mode=(stat.S_IFREG | 0o755), st_ctime=now, st_mtime=now, st_atime=now, st_nlink=1,
                    st_size=self.size)


class BlockchainFS(Operations):
    print "enters the Mini BlockChain FS"
    def __init__(self):
        self.fd = 0
        self.data = Custom()
        # print "print the entire data", self.data.data
        self.cache = {}
        self.cache['/'] = DirEntry()
        # print "cache", self.cache['/']
    # print "haha"

    def open(self, path, flags):
        self.fd += 1
        return self.fd

    # print "haha"

    def statfs(self, path):
        return dict(f_bsize=512, f_blocks=4096, f_bavail=2048)

    # print "haha"

    def getattr(self, path, fh=None):
        # print "enters getattr", path
        if path in self.cache:
            return self.cache[path].stat()
        raise FuseOSError(errno.ENOENT)

    def read(self, path, size, offset, fh):
        # print"baby pathndfjdsjfjdfjdf", path
        if path in self.cache:
            # print "hi", self.cache[path].data[offset:offset + size]
            return self.cache[path].data[offset:offset + size]
        return '\x00' * size

    def readdir(self, path, fh):
        # print "hahagsgdkasd"
        if path == '/':  # /
            l = ['%03dxxx' % x for x in range(0, 10 + 1)]
            for i in l:
                self.cache[path + i] = DirEntry()
            return ['.', '..'] + l

        if re.match(r'^/[0-9]{3}xxx$', path):  # /000xxx
            offset = int(path[1:4])
            l = ['%06d' % x for x in range(0, min((offset + 1) * 1000, 10 + 1))]
            for i in l:
                self.cache[path + '/' + i] = DirEntry()
            return ['.', '..'] + l

        if re.match(r'^/[0-9]{3}xxx/[0-9]{6}$', path):  # /000xxx/000000
            # print "match found", int(path[8:14])
            blockhash = self.data.blockhash_by_index(int(path[8:14]))
            # print "hi hash", blockhash
            blockinfo = self.data.blockinfo(blockhash)
            # print "bin", blockinfo
            toadd = {
                    'blockhash': FileEntry(blockhash),
                    'nonce': FileEntry(str(blockinfo['Nonce'])),
                    'previousblockhash': FileEntry(str(blockinfo['PrevHash'])),
                    'transactions': FileEntry(str(len(blockinfo['Transactions']))),
                    'miner': FileEntry(str(blockinfo['MinerID']))
            }
            for tx in blockinfo['Transactions'][0]:
                # print "tx here is ", tx
                # print "transa", blockinfo['Transactions'][0]
                if tx == "UUID":
                    toadd[blockinfo['Transactions'][0]['UUID']] = DirEntry()
            for i in toadd:
                self.cache[path + '/' + i] = toadd[i]
            return ['.', '..'] + toadd.keys()
        # print "baby check the path here", path
        if re.match(r'^/[0-9]{3}xxx/[0-9]{6}/[0-9a-f]{32}$', path):
            # print "now txinfo f3 with path", path  # /000xxx
            # /000000/ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff
            txinfo = self.data.txinfo(path[-32:])
            toadd = {
                'MiningFee': FileEntry(str(txinfo['MiningFee'])),
                'Value': FileEntry(str(txinfo['Value'])),
                'ToID': FileEntry(str(len(txinfo['ToID']))),
                'FromID': FileEntry(str(txinfo['FromID'])),
                'Type': FileEntry(str(len(txinfo['Type'])))
            }
            for i in toadd:
                self.cache[path + '/' + i] = toadd[i]
                # print "noclue", self.cache[path + '/' + i]
            return ['.', '..'] + toadd.keys()

        return ['.', '..']


def main():
    if len(sys.argv) != 2:
        print('usage: %s <mountpoint>' % sys.argv[0])
        sys.exit(1)
    print("Welcome to Mini BlockChain Filesystem!!!")
    fuse = FUSE(BlockchainFS(), sys.argv[1], foreground=True, **{'allow_other': True})


if __name__ == '__main__':
    main()
