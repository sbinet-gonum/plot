// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tex

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func unicodeIndex(v string, math bool) rune {
	if !math {
		r, _ := utf8.DecodeRune([]byte(v))
		if r == utf8.RuneError {
			panic(fmt.Errorf("tex: invalid rune %q", v))
		}
		return r
	}
	// From UTF #25: U+2212 minus sign is the preferred
	// representation of the unary and binary minus sign rather than
	// the ASCII-derived U+002D hyphen-minus, because minus sign is
	// unambiguous and because it is rendered with a more desirable
	// length, usually longer than a hyphen.
	if v == "-" {
		return 0x2212
	}

	r, _ := utf8.DecodeRune([]byte(v))
	if r != utf8.RuneError {
		return r
	}

	r, ok := tex2uni[strings.Replace(v, `\`, "", 1)]
	if ok {
		return r
	}

	panic(fmt.Errorf("%q is not a valid unicode character nor a known TeX symbol", v))
}

var (
	tex2uni = map[string]rune{
		`widehat`:                  0x0302,
		`widetilde`:                0x0303,
		`widebar`:                  0x0305,
		`langle`:                   0x27e8,
		`rangle`:                   0x27e9,
		`perp`:                     0x27c2,
		`neq`:                      0x2260,
		`Join`:                     0x2a1d,
		`leqslant`:                 0x2a7d,
		`geqslant`:                 0x2a7e,
		`lessapprox`:               0x2a85,
		`gtrapprox`:                0x2a86,
		`lesseqqgtr`:               0x2a8b,
		`gtreqqless`:               0x2a8c,
		`triangleeq`:               0x225c,
		`eqslantless`:              0x2a95,
		`eqslantgtr`:               0x2a96,
		`backepsilon`:              0x03f6,
		`precapprox`:               0x2ab7,
		`succapprox`:               0x2ab8,
		`fallingdotseq`:            0x2252,
		`subseteqq`:                0x2ac5,
		`supseteqq`:                0x2ac6,
		`varpropto`:                0x221d,
		`precnapprox`:              0x2ab9,
		`succnapprox`:              0x2aba,
		`subsetneqq`:               0x2acb,
		`supsetneqq`:               0x2acc,
		`lnapprox`:                 0x2ab9,
		`gnapprox`:                 0x2aba,
		`longleftarrow`:            0x27f5,
		`longrightarrow`:           0x27f6,
		`longleftrightarrow`:       0x27f7,
		`Longleftarrow`:            0x27f8,
		`Longrightarrow`:           0x27f9,
		`Longleftrightarrow`:       0x27fa,
		`longmapsto`:               0x27fc,
		`leadsto`:                  0x21dd,
		`dashleftarrow`:            0x290e,
		`dashrightarrow`:           0x290f,
		`circlearrowleft`:          0x21ba,
		`circlearrowright`:         0x21bb,
		`leftrightsquigarrow`:      0x21ad,
		`leftsquigarrow`:           0x219c,
		`rightsquigarrow`:          0x219d,
		`Game`:                     0x2141,
		`hbar`:                     0x0127,
		`hslash`:                   0x210f,
		`ldots`:                    0x2026,
		`vdots`:                    0x22ee,
		`doteqdot`:                 0x2251,
		`doteq`:                    8784,
		`partial`:                  8706,
		`gg`:                       8811,
		`asymp`:                    8781,
		`blacktriangledown`:        9662,
		`otimes`:                   8855,
		`nearrow`:                  8599,
		`varpi`:                    982,
		`vee`:                      8744,
		`vec`:                      8407,
		`smile`:                    8995,
		`succnsim`:                 8937,
		`gimel`:                    8503,
		`vert`:                     124,
		`|`:                        124,
		`varrho`:                   1009,
		`P`:                        182,
		`approxident`:              8779,
		`Swarrow`:                  8665,
		`textasciicircum`:          94,
		`imageof`:                  8887,
		`ntriangleleft`:            8938,
		`nleq`:                     8816,
		`div`:                      247,
		`nparallel`:                8742,
		`Leftarrow`:                8656,
		`lll`:                      8920,
		`oiint`:                    8751,
		`ngeq`:                     8817,
		`Theta`:                    920,
		`origof`:                   8886,
		`blacksquare`:              9632,
		`solbar`:                   9023,
		`neg`:                      172,
		`sum`:                      8721,
		`Vdash`:                    8873,
		`coloneq`:                  8788,
		`degree`:                   176,
		`bowtie`:                   8904,
		`blacktriangleright`:       9654,
		`varsigma`:                 962,
		`leq`:                      8804,
		`ggg`:                      8921,
		`lneqq`:                    8808,
		`scurel`:                   8881,
		`stareq`:                   8795,
		`BbbN`:                     8469,
		`nLeftarrow`:               8653,
		`nLeftrightarrow`:          8654,
		`k`:                        808,
		`bot`:                      8869,
		`BbbC`:                     8450,
		`Lsh`:                      8624,
		`leftleftarrows`:           8647,
		`BbbZ`:                     8484,
		`digamma`:                  989,
		`BbbR`:                     8477,
		`BbbP`:                     8473,
		`BbbQ`:                     8474,
		`vartriangleright`:         8883,
		`succsim`:                  8831,
		`wedge`:                    8743,
		`lessgtr`:                  8822,
		`veebar`:                   8891,
		`mapsdown`:                 8615,
		`Rsh`:                      8625,
		`chi`:                      967,
		`prec`:                     8826,
		`nsubseteq`:                8840,
		`therefore`:                8756,
		`eqcirc`:                   8790,
		`textexclamdown`:           161,
		`nRightarrow`:              8655,
		`flat`:                     9837,
		`notin`:                    8713,
		`llcorner`:                 8990,
		`varepsilon`:               949,
		`bigtriangleup`:            9651,
		`aleph`:                    8501,
		`dotminus`:                 8760,
		`upsilon`:                  965,
		`Lambda`:                   923,
		`cap`:                      8745,
		`barleftarrow`:             8676,
		`mu`:                       956,
		`boxplus`:                  8862,
		`mp`:                       8723,
		`circledast`:               8859,
		`tau`:                      964,
		`in`:                       8712,
		`backslash`:                92,
		`varnothing`:               8709,
		`sharp`:                    9839,
		`eqsim`:                    8770,
		`gnsim`:                    8935,
		`Searrow`:                  8664,
		`updownarrows`:             8645,
		`heartsuit`:                9825,
		`trianglelefteq`:           8884,
		`ddag`:                     8225,
		`sqsubseteq`:               8849,
		`mapsfrom`:                 8612,
		`boxbar`:                   9707,
		`sim`:                      8764,
		`Nwarrow`:                  8662,
		`nequiv`:                   8802,
		`succ`:                     8827,
		`vdash`:                    8866,
		`Leftrightarrow`:           8660,
		`parallel`:                 8741,
		`invnot`:                   8976,
		`natural`:                  9838,
		`ss`:                       223,
		`uparrow`:                  8593,
		`nsim`:                     8769,
		`hookrightarrow`:           8618,
		`Equiv`:                    8803,
		`approx`:                   8776,
		`Vvdash`:                   8874,
		`nsucc`:                    8833,
		`leftrightharpoons`:        8651,
		`Re`:                       8476,
		`boxminus`:                 8863,
		`equiv`:                    8801,
		`Lleftarrow`:               8666,
		`ll`:                       8810,
		`Cup`:                      8915,
		`measeq`:                   8798,
		`upharpoonleft`:            8639,
		`lq`:                       8216,
		`Upsilon`:                  933,
		`subsetneq`:                8842,
		`greater`:                  62,
		`supsetneq`:                8843,
		`Cap`:                      8914,
		`L`:                        321,
		`spadesuit`:                9824,
		`lrcorner`:                 8991,
		`not`:                      824,
		`bar`:                      772,
		`rightharpoonaccent`:       8401,
		`boxdot`:                   8865,
		`l`:                        322,
		`leftharpoondown`:          8637,
		`bigcup`:                   8899,
		`iint`:                     8748,
		`bigwedge`:                 8896,
		`downharpoonleft`:          8643,
		`textasciitilde`:           126,
		`subset`:                   8834,
		`leqq`:                     8806,
		`mapsup`:                   8613,
		`nvDash`:                   8877,
		`looparrowleft`:            8619,
		`nless`:                    8814,
		`rightarrowbar`:            8677,
		`Vert`:                     8214,
		`downdownarrows`:           8650,
		`uplus`:                    8846,
		`simeq`:                    8771,
		`napprox`:                  8777,
		`ast`:                      8727,
		`twoheaduparrow`:           8607,
		`doublebarwedge`:           8966,
		`Sigma`:                    931,
		`leftharpoonaccent`:        8400,
		`ntrianglelefteq`:          8940,
		`nexists`:                  8708,
		`times`:                    215,
		`measuredangle`:            8737,
		`bumpeq`:                   8783,
		`carriagereturn`:           8629,
		`adots`:                    8944,
		`checkmark`:                10003,
		`lambda`:                   955,
		`xi`:                       958,
		`rbrace`:                   125,
		`rbrack`:                   93,
		`Nearrow`:                  8663,
		`maltese`:                  10016,
		`clubsuit`:                 9827,
		`top`:                      8868,
		`overarc`:                  785,
		`varphi`:                   966,
		`Delta`:                    916,
		`iota`:                     953,
		`nleftarrow`:               8602,
		`candra`:                   784,
		`supset`:                   8835,
		`triangleleft`:             9665,
		`gtreqless`:                8923,
		`ntrianglerighteq`:         8941,
		`quad`:                     8195,
		`Xi`:                       926,
		`gtrdot`:                   8919,
		`leftthreetimes`:           8907,
		`minus`:                    8722,
		`preccurlyeq`:              8828,
		`nleftrightarrow`:          8622,
		`lambdabar`:                411,
		`blacktriangle`:            9652,
		`kernelcontraction`:        8763,
		`Phi`:                      934,
		`angle`:                    8736,
		`spadesuitopen`:            9828,
		`eqless`:                   8924,
		`mid`:                      8739,
		`varkappa`:                 1008,
		`Ldsh`:                     8626,
		`updownarrow`:              8597,
		`beta`:                     946,
		`textquotedblleft`:         8220,
		`rho`:                      961,
		`alpha`:                    945,
		`intercal`:                 8890,
		`beth`:                     8502,
		`grave`:                    768,
		`acwopencirclearrow`:       8634,
		`nmid`:                     8740,
		`nsupset`:                  8837,
		`sigma`:                    963,
		`dot`:                      775,
		`Rightarrow`:               8658,
		`turnednot`:                8985,
		`backsimeq`:                8909,
		`leftarrowtail`:            8610,
		`approxeq`:                 8778,
		`curlyeqsucc`:              8927,
		`rightarrowtail`:           8611,
		`Psi`:                      936,
		`copyright`:                169,
		`yen`:                      165,
		`vartriangleleft`:          8882,
		`rasp`:                     700,
		`triangleright`:            9655,
		`precsim`:                  8830,
		`infty`:                    8734,
		`geq`:                      8805,
		`updownarrowbar`:           8616,
		`precnsim`:                 8936,
		`H`:                        779,
		`ulcorner`:                 8988,
		`looparrowright`:           8620,
		`ncong`:                    8775,
		`downarrow`:                8595,
		`circeq`:                   8791,
		`subseteq`:                 8838,
		`bigstar`:                  9733,
		`prime`:                    8242,
		`lceil`:                    8968,
		`Rrightarrow`:              8667,
		`oiiint`:                   8752,
		`curlywedge`:               8911,
		`vDash`:                    8872,
		`lfloor`:                   8970,
		`ddots`:                    8945,
		`exists`:                   8707,
		`underbar`:                 817,
		`Pi`:                       928,
		`leftrightarrows`:          8646,
		`sphericalangle`:           8738,
		`coprod`:                   8720,
		`circledcirc`:              8858,
		`gtrsim`:                   8819,
		`gneqq`:                    8809,
		`between`:                  8812,
		`theta`:                    952,
		`complement`:               8705,
		`arceq`:                    8792,
		`nVdash`:                   8878,
		`S`:                        167,
		`wr`:                       8768,
		`wp`:                       8472,
		`backcong`:                 8780,
		`lasp`:                     701,
		`c`:                        807,
		`nabla`:                    8711,
		`dotplus`:                  8724,
		`eta`:                      951,
		`forall`:                   8704,
		`eth`:                      240,
		`colon`:                    58,
		`sqcup`:                    8852,
		`rightrightarrows`:         8649,
		`sqsupset`:                 8848,
		`mapsto`:                   8614,
		`bigtriangledown`:          9661,
		`sqsupseteq`:               8850,
		`propto`:                   8733,
		`pi`:                       960,
		`pm`:                       177,
		`dots`:                     0x2026,
		`nrightarrow`:              8603,
		`textasciiacute`:           180,
		`Doteq`:                    8785,
		`breve`:                    774,
		`sqcap`:                    8851,
		`twoheadrightarrow`:        8608,
		`kappa`:                    954,
		`vartriangle`:              9653,
		`diamondsuit`:              9826,
		`pitchfork`:                8916,
		`blacktriangleleft`:        9664,
		`nprec`:                    8832,
		`curvearrowright`:          8631,
		`barwedge`:                 8892,
		`multimap`:                 8888,
		`textquestiondown`:         191,
		`cong`:                     8773,
		`rtimes`:                   8906,
		`rightzigzagarrow`:         8669,
		`rightarrow`:               8594,
		`leftarrow`:                8592,
		`__sqrt__`:                 8730,
		`twoheaddownarrow`:         8609,
		`oint`:                     8750,
		`bigvee`:                   8897,
		`eqdef`:                    8797,
		`sterling`:                 163,
		`phi`:                      981,
		`Updownarrow`:              8661,
		`backprime`:                8245,
		`emdash`:                   8212,
		`Gamma`:                    915,
		`i`:                        305,
		`rceil`:                    8969,
		`leftharpoonup`:            8636,
		`Im`:                       8465,
		`curvearrowleft`:           8630,
		`wedgeq`:                   8793,
		`curlyeqprec`:              8926,
		`questeq`:                  8799,
		`less`:                     60,
		`upuparrows`:               8648,
		`tilde`:                    771,
		`textasciigrave`:           96,
		`smallsetminus`:            8726,
		`ell`:                      8467,
		`cup`:                      8746,
		`danger`:                   9761,
		`nVDash`:                   8879,
		`cdotp`:                    183,
		`cdots`:                    8943,
		`hat`:                      770,
		`eqgtr`:                    8925,
		`psi`:                      968,
		`frown`:                    8994,
		`acute`:                    769,
		`downzigzagarrow`:          8623,
		`ntriangleright`:           8939,
		`cupdot`:                   8845,
		`circleddash`:              8861,
		`oslash`:                   8856,
		`mho`:                      8487,
		`d`:                        803,
		`sqsubset`:                 8847,
		`cdot`:                     8901,
		`Omega`:                    937,
		`OE`:                       338,
		`veeeq`:                    8794,
		`Finv`:                     8498,
		`t`:                        865,
		`leftrightarrow`:           8596,
		`swarrow`:                  8601,
		`rightthreetimes`:          8908,
		`rightleftharpoons`:        8652,
		`lesssim`:                  8818,
		`searrow`:                  8600,
		`because`:                  8757,
		`gtrless`:                  8823,
		`star`:                     8902,
		`nsubset`:                  8836,
		`zeta`:                     950,
		`dddot`:                    8411,
		`bigcirc`:                  9675,
		`Supset`:                   8913,
		`circ`:                     8728,
		`slash`:                    8725,
		`ocirc`:                    778,
		`prod`:                     8719,
		`twoheadleftarrow`:         8606,
		`daleth`:                   8504,
		`upharpoonright`:           8638,
		`odot`:                     8857,
		`Uparrow`:                  8657,
		`O`:                        216,
		`hookleftarrow`:            8617,
		`trianglerighteq`:          8885,
		`nsime`:                    8772,
		`oe`:                       339,
		`nwarrow`:                  8598,
		`o`:                        248,
		`ddddot`:                   8412,
		`downharpoonright`:         8642,
		`succcurlyeq`:              8829,
		`gamma`:                    947,
		`scrR`:                     8475,
		`dag`:                      8224,
		`thickspace`:               8197,
		`frakZ`:                    8488,
		`lessdot`:                  8918,
		`triangledown`:             9663,
		`ltimes`:                   8905,
		`scrB`:                     8492,
		`endash`:                   8211,
		`scrE`:                     8496,
		`scrF`:                     8497,
		`scrH`:                     8459,
		`scrI`:                     8464,
		`rightharpoondown`:         8641,
		`scrL`:                     8466,
		`scrM`:                     8499,
		`frakC`:                    8493,
		`nsupseteq`:                8841,
		`circledR`:                 174,
		`circledS`:                 9416,
		`ngtr`:                     8815,
		`bigcap`:                   8898,
		`scre`:                     8495,
		`Downarrow`:                8659,
		`scrg`:                     8458,
		`overleftrightarrow`:       8417,
		`scro`:                     8500,
		`lnsim`:                    8934,
		`eqcolon`:                  8789,
		`curlyvee`:                 8910,
		`urcorner`:                 8989,
		`lbrace`:                   123,
		`Bumpeq`:                   8782,
		`delta`:                    948,
		`boxtimes`:                 8864,
		`overleftarrow`:            8406,
		`prurel`:                   8880,
		`clubsuitopen`:             9831,
		`cwopencirclearrow`:        8635,
		`geqq`:                     8807,
		`rightleftarrows`:          8644,
		`ac`:                       8766,
		`ae`:                       230,
		`int`:                      8747,
		`rfloor`:                   8971,
		`risingdotseq`:             8787,
		`nvdash`:                   8876,
		`diamond`:                  8900,
		`ddot`:                     776,
		`backsim`:                  8765,
		`oplus`:                    8853,
		`triangleq`:                8796,
		`check`:                    780,
		`ni`:                       8715,
		`iiint`:                    8749,
		`ne`:                       8800,
		`lesseqgtr`:                8922,
		`obar`:                     9021,
		`supseteq`:                 8839,
		`nu`:                       957,
		`AA`:                       197,
		`AE`:                       198,
		`models`:                   8871,
		`ominus`:                   8854,
		`dashv`:                    8867,
		`omega`:                    969,
		`rq`:                       8217,
		`Subset`:                   8912,
		`rightharpoonup`:           8640,
		`Rdsh`:                     8627,
		`bullet`:                   8729,
		`divideontimes`:            8903,
		`lbrack`:                   91,
		`textquotedblright`:        8221,
		`Colon`:                    8759,
		`%`:                        37,
		`$`:                        36,
		`{`:                        123,
		`}`:                        125,
		`_`:                        95,
		`#`:                        35,
		`imath`:                    0x131,
		`circumflexaccent`:         770,
		`combiningbreve`:           774,
		`combiningoverline`:        772,
		`combininggraveaccent`:     768,
		`combiningacuteaccent`:     769,
		`combiningdiaeresis`:       776,
		`combiningtilde`:           771,
		`combiningrightarrowabove`: 8407,
		`combiningdotabove`:        775,
		`to`:                       8594,
		`succeq`:                   8829,
		`emptyset`:                 8709,
		`leftparen`:                40,
		`rightparen`:               41,
		`bigoplus`:                 10753,
		`leftangle`:                10216,
		`rightangle`:               10217,
		`leftbrace`:                124,
		`rightbrace`:               125,
		`jmath`:                    567,
		`bigodot`:                  10752,
		`preceq`:                   8828,
		`biguplus`:                 10756,
		`epsilon`:                  949,
		`vartheta`:                 977,
		`bigotimes`:                10754,
		`guillemotleft`:            171,
		`ring`:                     730,
		`Thorn`:                    222,
		`guilsinglright`:           8250,
		`perthousand`:              8240,
		`macron`:                   175,
		`cent`:                     162,
		`guillemotright`:           187,
		`equal`:                    61,
		`asterisk`:                 42,
		`guilsinglleft`:            8249,
		`plus`:                     43,
		`thorn`:                    254,
		`dagger`:                   8224,
	}
)
