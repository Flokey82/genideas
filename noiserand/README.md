# noiserand

The `noiserand` package implements a random number generator that is based on multidimensional noise, using one of the dimensions as a seed, which allows slight alterations to the seed to slightly alter the generated sequence.

## What is this good for?

Well, for example it can be used for procedural generation and if you have found a seed that produces an interesting result, you can explore neighboring seeds which will produce slight variations of the original output.

## Usage

To use the `noiserand` package, first import it into your Go code:

```go
import "github.com/Flokey82/genideas/noiserand"
```

Then, create a new `NoiseRand` object with a random seed and a noise function:

```go
r := noiserand.New(12345.0, func(x, y float64) float64 {
    // Replace this with your own noise function.
    return 0.0
})
```

You can then use the `Intn`, `Int63`, and `Float64` methods to generate random numbers:

```go
// Generate a random integer in the range [0, 10).
n := r.Intn(10)

// Generate a random 63-bit integer.
n64 := r.Int63()

// Generate a random float in the range [0.0, 1.0).
f := r.Float64()

// Generate a pseudo-random permutation of the integers [0,n).
p := r.Perm(10)
```

## Contributions

We welcome contributions to the noiserand package! If you have ideas for improvements or new features, please feel free to submit a pull request on GitHub. We also appreciate bug reports and feedback on how the package is being used. By working together, we can make the noiserand package even better for everyone. Thank you for your support!

Once this package is in a good state, I will move it to its own repository.

## License

The `noiserand` package is licensed under the MIT License. See the LICENSE file for more information.

