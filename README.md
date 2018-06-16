## Research Notes

- ed25519 curve and related vrf seem to have been proposed for x/crypto here: https://go-review.googlesource.com/c/crypto/+/13798
- since then (?) curve25519 and ed25519 have made it into x/crypto: https://godoc.org/golang.org/x/crypto
- vrf seems to have been borrowed (verbatim?) in the coniks-go project: https://github.com/coniks-sys/coniks-go/tree/master/crypto/vrf

### Implementation Notes

- Sortition
    - *Had* some concerns about precision in the float64 calculation of probabilities in `getSortitionIntervals`.
    Spent some time investigating a large-integer approach. This approach has some serious performance issues because the exponents get *huge* very quickly.
    Also showed that in reasonable scenarios it appears that the probability 'mass' that is lost is ~10^16 which is obviously negligible.
    Seems I might have 'discovered' the lower-bound precision of float64, which is `1.11 × 10−16`: https://en.wikipedia.org/wiki/Double-precision_floating-point_format
    - If the interval calculations are left in terms of probabilities (likely *are*), it may be worth circling back on some performance enhancements in `getSortitionIntervals`.
    As noted in comments there, my initial implementation following the Algorand paper directly but left some fairly simple optimizations out.
    However, this also may not end up being significant if the interval caching mechanism ends up being useful.
    - Also worth noting for that a given tau, total weights and user weight, *all* users with that weight will have the same intervals.
    The current 'caching' mechanism is per-user but this implies that a more 'global' mechanism may have benefits too.
    - I haven't really spent any time at all reviewing a/o testing the VRF implementation or the elliptic curve implementation it relies on.
    Might be desirable (necessary?) to do so in the future.