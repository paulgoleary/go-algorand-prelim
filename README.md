## Research Notes

- ed25519 curve and related vrf seem to have been proposed for x/crypto here: https://go-review.googlesource.com/c/crypto/+/13798
- since then (?) curve25519 and ed25519 have made it into x/crypto: https://godoc.org/golang.org/x/crypto
- vrf seems to have been borrowed (verbatim?) in the coniks-go project: https://github.com/coniks-sys/coniks-go/tree/master/crypto/vrf

### Implementation Notes

- Sortition
    - I have some concerns about the precision limits in the current implementation. Specifically, I could not readily find any large-float support for exponentiation in golang.
    As exponentiation is required for the probability mass function - p^k(1-p)^n-k - this might cause some significant loss of precision, particularly when the total weight of tokens gets large which is the denominator of the probability.
    - Spent some time dinking around with the math to see if I could do a comparable calculation with all large-integer math. Seems plausible.
    The thought here is that the intervals do *not* have to end up being expressed as probabilities. The 'search' to find `j` is with the ratio of the VRF hash with it's maximum possible integer value.
    If the buckets could be expressed in terms of integer ranges that cover the same maximum integer value it might be possible to keep everything from losing precision.
    It may seem like this would be slower but we're already doing many calculations with large integers and floats, not to mention VRF which does large-field elliptic curve calculations.
    - If the interval calculations are left in terms of probabilities, it may be worth circling back on some performance enhancements in `getSortitionIntervals`.
    As noted in comments there, my initial implementation following the Algorand paper directly but left some fairly simple optimizations out.
    However, this also may not end up being significant if the interval caching mechanism ends up being useful.
    - Also worth noting for that a given tau, total weights and user weight, *all* users with that weight will have the same intervals.
    The current 'caching' mechanism is per-user but this implies that a more 'global' mechanism may have benefits too.
    - I haven't really spent any time at all reviewing a/o testing the VRF implementation or the elliptic curve implementation it relies on.
    Might be desirable (necessary?) to do so in the future.