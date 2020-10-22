# Bloom Filter Password Checker

An implementing of [NIST 800-63-3b Leaked Password Check](https://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-63b.pdf) using a [blooming filter](https://dl.acm.org/citation.cfm?doid=362686.362692) built from the [Have I been pwned](https://haveibeenpwned.com/Passwords) SHA1 password hash list. The Have I Been Pwned list contains more than 580 million password hashes and is 25GB uncompressed (as of Jun 2020). A bloom filter of this list is about 900MB and will fit entirely into memory on a virtual machine or Docker container with just 2GB of RAM.

## Why a Bloom Filter?

It's one of the simplest, smallest and fastest data structures for this task. Bloom filters have constant time O(1) performance (where K is the constant) for insertion and lookup. K is the number of times a password is hashed. Bloom filters can easily handle billions of banned password hashes with very modest resources.

## Partial SHA1 Hashes

SHA1 hashes are 20 bytes of raw binary data and thus typically hex encoded for a total of 40 characters. Blooming Password uses just the first 16 hex encoded characters of the hashes to build the bloom filter and to test the filter for membership. The program rejects complete hashes if they are sent. False positive rates in the bloom filter are not impacted by the shortening of the SHA1 password hashes. The cardinality of the set is unchanged. The false positive rate is 1 out of 1000. You may verify the cardinality is unchanged after truncating the hashes.
