package main_test

import (
	"net/http"
	"strings"
	"testing"
	"worldcoin/gnark-mbu/logging"
	"worldcoin/gnark-mbu/prover"
	"worldcoin/gnark-mbu/server"

	gnarkLogger "github.com/consensys/gnark/logger"
)

const ProverAddress = "localhost:8080"
const MetricsAddress = "localhost:9999"

var mode string

func TestMain(m *testing.M) {
	gnarkLogger.Set(*logging.Logger())
	logging.Logger().Info().Msg("Setting up the prover")
	ps, err := prover.SetupInsertion(22, 3)
	if err != nil {
		panic(err)
	}
	cfg := server.Config{
		ProverAddress:  ProverAddress,
		MetricsAddress: MetricsAddress,
		Mode:           server.InsertionMode,
	}
	logging.Logger().Info().Msg("Starting the insertion server")
	instance := server.Run(&cfg, ps)
	logging.Logger().Info().Msg("Running the insertion tests")
	mode = server.InsertionMode
	m.Run()
	instance.RequestStop()
	instance.AwaitStop()
	// cfg.Mode = server.DeletionMode
	// ps, err = prover.SetupDeletion(22, 1)
	// if err != nil {
	// 	panic(err)
	// }
	// logging.Logger().Info().Msg("Starting the deletion server")
	// instance = server.Run(&cfg, ps)
	// logging.Logger().Info().Msg("Running the deletion tests")
	// mode = server.DeletionMode
	// m.Run()
	// instance.RequestStop()
	// instance.AwaitStop()
}

// func TestWrongMethod(t *testing.T) {
// 	response, err := http.Get("http://localhost:8080/prove")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if response.StatusCode != http.StatusMethodNotAllowed {
// 		t.Fatalf("Expected status code %d, got %d", http.StatusMethodNotAllowed, response.StatusCode)
// 	}
// }

// func TestInsertionHappyPath22_1(t *testing.T) {
// 	if mode != server.InsertionMode {
// 		return
// 	}
// 	body := `{
// 		"inputHash":"0x9f565aa05c660b0b46c446d95ad75716532764a42cdbbe5f4a9bbc7e01a1f624",
// 		"startIndex":0,
// 		"preRoot":"0x18cca2a66b5c0787981e69aefd84852d74af0e93ef4912b4648c05f722efe52b",
// 		"postRoot":"0xb15b7b842cd0270ec5b994b28e9a62428c9cad63f35f83a82672acd33fbca9a",
// 		"identityCommitments":["0x1"],
// 		"merkleProofs":[
// 			["0x0",
// 			"0x2098f5fb9e239eab3ceac3f27b81e481dc3124d55ffed523a839ee8446b64864",
// 			"0x1069673dcdb12263df301a6ff584a7ec261a44cb9dc68df067a4774460b1f1e1",
// 			"0x18f43331537ee2af2e3d758d50f72106467c6eea50371dd528d57eb2b856d238",
// 			"0x7f9d837cb17b0d36320ffe93ba52345f1b728571a568265caac97559dbc952a",
// 			"0x2b94cf5e8746b3f5c9631f4c5df32907a699c58c94b2ad4d7b5cec1639183f55",
// 			"0x2dee93c5a666459646ea7d22cca9e1bcfed71e6951b953611d11dda32ea09d78",
// 			"0x78295e5a22b84e982cf601eb639597b8b0515a88cb5ac7fa8a4aabe3c87349d",
// 			"0x2fa5e5f18f6027a6501bec864564472a616b2e274a41211a444cbe3a99f3cc61",
// 			"0xe884376d0d8fd21ecb780389e941f66e45e7acce3e228ab3e2156a614fcd747",
// 			"0x1b7201da72494f1e28717ad1a52eb469f95892f957713533de6175e5da190af2",
// 			"0x1f8d8822725e36385200c0b201249819a6e6e1e4650808b5bebc6bface7d7636",
// 			"0x2c5d82f66c914bafb9701589ba8cfcfb6162b0a12acf88a8d0879a0471b5f85a",
// 			"0x14c54148a0940bb820957f5adf3fa1134ef5c4aaa113f4646458f270e0bfbfd0",
// 			"0x190d33b12f986f961e10c0ee44d8b9af11be25588cad89d416118e4bf4ebe80c",
// 			"0x22f98aa9ce704152ac17354914ad73ed1167ae6596af510aa5b3649325e06c92",
// 			"0x2a7c7c9b6ce5880b9f6f228d72bf6a575a526f29c66ecceef8b753d38bba7323",
// 			"0x2e8186e558698ec1c67af9c14d463ffc470043c9c2988b954d75dd643f36b992",
// 			"0xf57c5571e9a4eab49e2c8cf050dae948aef6ead647392273546249d1c1ff10f",
// 			"0x1830ee67b5fb554ad5f63d4388800e1cfe78e310697d46e43c9ce36134f72cca",
// 			"0x2134e76ac5d21aab186c2be1dd8f84ee880a1e46eaf712f9d371b6df22191f3e",
// 			"0x19df90ec844ebc4ffeebd866f33859b0c051d8c958ee3aa88f8f8df3db91a5b1"]
// 		]}`
// 	response, err := http.Post("http://localhost:8080/prove", "application/json", strings.NewReader(body))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if response.StatusCode != http.StatusOK {
// 		t.Fatalf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
// 	}
// }

// func TestInsertionHappyPath22_2(t *testing.T) {
// 	if mode != server.InsertionMode {
// 		return
// 	}
// 	body := `{"inputHash":"0x4e9204bee0007caccbd3278dab3aa3211b58726962d98405f28cf24af5f88331","startIndex":0,"preRoot":"0x18cca2a66b5c0787981e69aefd84852d74af0e93ef4912b4648c05f722efe52b","postRoot":"0x244b5b2e2a2cfe267cbc21dbb79689962f2b3a3234f517ff74972876480d5fb1","identityCommitments":["0x1","0x2"],"merkleProofs":[["0x0","0x2098f5fb9e239eab3ceac3f27b81e481dc3124d55ffed523a839ee8446b64864","0x1069673dcdb12263df301a6ff584a7ec261a44cb9dc68df067a4774460b1f1e1","0x18f43331537ee2af2e3d758d50f72106467c6eea50371dd528d57eb2b856d238","0x7f9d837cb17b0d36320ffe93ba52345f1b728571a568265caac97559dbc952a","0x2b94cf5e8746b3f5c9631f4c5df32907a699c58c94b2ad4d7b5cec1639183f55","0x2dee93c5a666459646ea7d22cca9e1bcfed71e6951b953611d11dda32ea09d78","0x78295e5a22b84e982cf601eb639597b8b0515a88cb5ac7fa8a4aabe3c87349d","0x2fa5e5f18f6027a6501bec864564472a616b2e274a41211a444cbe3a99f3cc61","0xe884376d0d8fd21ecb780389e941f66e45e7acce3e228ab3e2156a614fcd747","0x1b7201da72494f1e28717ad1a52eb469f95892f957713533de6175e5da190af2","0x1f8d8822725e36385200c0b201249819a6e6e1e4650808b5bebc6bface7d7636","0x2c5d82f66c914bafb9701589ba8cfcfb6162b0a12acf88a8d0879a0471b5f85a","0x14c54148a0940bb820957f5adf3fa1134ef5c4aaa113f4646458f270e0bfbfd0","0x190d33b12f986f961e10c0ee44d8b9af11be25588cad89d416118e4bf4ebe80c","0x22f98aa9ce704152ac17354914ad73ed1167ae6596af510aa5b3649325e06c92","0x2a7c7c9b6ce5880b9f6f228d72bf6a575a526f29c66ecceef8b753d38bba7323","0x2e8186e558698ec1c67af9c14d463ffc470043c9c2988b954d75dd643f36b992","0xf57c5571e9a4eab49e2c8cf050dae948aef6ead647392273546249d1c1ff10f","0x1830ee67b5fb554ad5f63d4388800e1cfe78e310697d46e43c9ce36134f72cca","0x2134e76ac5d21aab186c2be1dd8f84ee880a1e46eaf712f9d371b6df22191f3e","0x19df90ec844ebc4ffeebd866f33859b0c051d8c958ee3aa88f8f8df3db91a5b1"],["0x1","0x2098f5fb9e239eab3ceac3f27b81e481dc3124d55ffed523a839ee8446b64864","0x1069673dcdb12263df301a6ff584a7ec261a44cb9dc68df067a4774460b1f1e1","0x18f43331537ee2af2e3d758d50f72106467c6eea50371dd528d57eb2b856d238","0x7f9d837cb17b0d36320ffe93ba52345f1b728571a568265caac97559dbc952a","0x2b94cf5e8746b3f5c9631f4c5df32907a699c58c94b2ad4d7b5cec1639183f55","0x2dee93c5a666459646ea7d22cca9e1bcfed71e6951b953611d11dda32ea09d78","0x78295e5a22b84e982cf601eb639597b8b0515a88cb5ac7fa8a4aabe3c87349d","0x2fa5e5f18f6027a6501bec864564472a616b2e274a41211a444cbe3a99f3cc61","0xe884376d0d8fd21ecb780389e941f66e45e7acce3e228ab3e2156a614fcd747","0x1b7201da72494f1e28717ad1a52eb469f95892f957713533de6175e5da190af2","0x1f8d8822725e36385200c0b201249819a6e6e1e4650808b5bebc6bface7d7636","0x2c5d82f66c914bafb9701589ba8cfcfb6162b0a12acf88a8d0879a0471b5f85a","0x14c54148a0940bb820957f5adf3fa1134ef5c4aaa113f4646458f270e0bfbfd0","0x190d33b12f986f961e10c0ee44d8b9af11be25588cad89d416118e4bf4ebe80c","0x22f98aa9ce704152ac17354914ad73ed1167ae6596af510aa5b3649325e06c92","0x2a7c7c9b6ce5880b9f6f228d72bf6a575a526f29c66ecceef8b753d38bba7323","0x2e8186e558698ec1c67af9c14d463ffc470043c9c2988b954d75dd643f36b992","0xf57c5571e9a4eab49e2c8cf050dae948aef6ead647392273546249d1c1ff10f","0x1830ee67b5fb554ad5f63d4388800e1cfe78e310697d46e43c9ce36134f72cca","0x2134e76ac5d21aab186c2be1dd8f84ee880a1e46eaf712f9d371b6df22191f3e","0x19df90ec844ebc4ffeebd866f33859b0c051d8c958ee3aa88f8f8df3db91a5b1"]]}`
// 	response, err := http.Post("http://localhost:8080/prove", "application/json", strings.NewReader(body))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if response.StatusCode != http.StatusOK {
// 		t.Fatalf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
// 	}
// }

func TestInsertionHappyPath22_3(t *testing.T) {
	if mode != server.InsertionMode {
		return
	}
	body := `{"inputHash":"0x92a772c17b6d94d4c2e76e85e53eb5f373f1a55041ee253ceef7203bf5025c32","startIndex":0,"preRoot":"0x18cca2a66b5c0787981e69aefd84852d74af0e93ef4912b4648c05f722efe52b","postRoot":"0x212a83a477a0540e4321b319c22f3245b93b2dddc41de9924597abef4113f60","identityCommitments":["0x1","0x2","0x3"],"merkleProofs":[["0x0","0x2098f5fb9e239eab3ceac3f27b81e481dc3124d55ffed523a839ee8446b64864","0x1069673dcdb12263df301a6ff584a7ec261a44cb9dc68df067a4774460b1f1e1","0x18f43331537ee2af2e3d758d50f72106467c6eea50371dd528d57eb2b856d238","0x7f9d837cb17b0d36320ffe93ba52345f1b728571a568265caac97559dbc952a","0x2b94cf5e8746b3f5c9631f4c5df32907a699c58c94b2ad4d7b5cec1639183f55","0x2dee93c5a666459646ea7d22cca9e1bcfed71e6951b953611d11dda32ea09d78","0x78295e5a22b84e982cf601eb639597b8b0515a88cb5ac7fa8a4aabe3c87349d","0x2fa5e5f18f6027a6501bec864564472a616b2e274a41211a444cbe3a99f3cc61","0xe884376d0d8fd21ecb780389e941f66e45e7acce3e228ab3e2156a614fcd747","0x1b7201da72494f1e28717ad1a52eb469f95892f957713533de6175e5da190af2","0x1f8d8822725e36385200c0b201249819a6e6e1e4650808b5bebc6bface7d7636","0x2c5d82f66c914bafb9701589ba8cfcfb6162b0a12acf88a8d0879a0471b5f85a","0x14c54148a0940bb820957f5adf3fa1134ef5c4aaa113f4646458f270e0bfbfd0","0x190d33b12f986f961e10c0ee44d8b9af11be25588cad89d416118e4bf4ebe80c","0x22f98aa9ce704152ac17354914ad73ed1167ae6596af510aa5b3649325e06c92","0x2a7c7c9b6ce5880b9f6f228d72bf6a575a526f29c66ecceef8b753d38bba7323","0x2e8186e558698ec1c67af9c14d463ffc470043c9c2988b954d75dd643f36b992","0xf57c5571e9a4eab49e2c8cf050dae948aef6ead647392273546249d1c1ff10f","0x1830ee67b5fb554ad5f63d4388800e1cfe78e310697d46e43c9ce36134f72cca","0x2134e76ac5d21aab186c2be1dd8f84ee880a1e46eaf712f9d371b6df22191f3e","0x19df90ec844ebc4ffeebd866f33859b0c051d8c958ee3aa88f8f8df3db91a5b1"],["0x1","0x2098f5fb9e239eab3ceac3f27b81e481dc3124d55ffed523a839ee8446b64864","0x1069673dcdb12263df301a6ff584a7ec261a44cb9dc68df067a4774460b1f1e1","0x18f43331537ee2af2e3d758d50f72106467c6eea50371dd528d57eb2b856d238","0x7f9d837cb17b0d36320ffe93ba52345f1b728571a568265caac97559dbc952a","0x2b94cf5e8746b3f5c9631f4c5df32907a699c58c94b2ad4d7b5cec1639183f55","0x2dee93c5a666459646ea7d22cca9e1bcfed71e6951b953611d11dda32ea09d78","0x78295e5a22b84e982cf601eb639597b8b0515a88cb5ac7fa8a4aabe3c87349d","0x2fa5e5f18f6027a6501bec864564472a616b2e274a41211a444cbe3a99f3cc61","0xe884376d0d8fd21ecb780389e941f66e45e7acce3e228ab3e2156a614fcd747","0x1b7201da72494f1e28717ad1a52eb469f95892f957713533de6175e5da190af2","0x1f8d8822725e36385200c0b201249819a6e6e1e4650808b5bebc6bface7d7636","0x2c5d82f66c914bafb9701589ba8cfcfb6162b0a12acf88a8d0879a0471b5f85a","0x14c54148a0940bb820957f5adf3fa1134ef5c4aaa113f4646458f270e0bfbfd0","0x190d33b12f986f961e10c0ee44d8b9af11be25588cad89d416118e4bf4ebe80c","0x22f98aa9ce704152ac17354914ad73ed1167ae6596af510aa5b3649325e06c92","0x2a7c7c9b6ce5880b9f6f228d72bf6a575a526f29c66ecceef8b753d38bba7323","0x2e8186e558698ec1c67af9c14d463ffc470043c9c2988b954d75dd643f36b992","0xf57c5571e9a4eab49e2c8cf050dae948aef6ead647392273546249d1c1ff10f","0x1830ee67b5fb554ad5f63d4388800e1cfe78e310697d46e43c9ce36134f72cca","0x2134e76ac5d21aab186c2be1dd8f84ee880a1e46eaf712f9d371b6df22191f3e","0x19df90ec844ebc4ffeebd866f33859b0c051d8c958ee3aa88f8f8df3db91a5b1"],["0x0","0x115cc0f5e7d690413df64c6b9662e9cf2a3617f2743245519e19607a4417189a","0x1069673dcdb12263df301a6ff584a7ec261a44cb9dc68df067a4774460b1f1e1","0x18f43331537ee2af2e3d758d50f72106467c6eea50371dd528d57eb2b856d238","0x7f9d837cb17b0d36320ffe93ba52345f1b728571a568265caac97559dbc952a","0x2b94cf5e8746b3f5c9631f4c5df32907a699c58c94b2ad4d7b5cec1639183f55","0x2dee93c5a666459646ea7d22cca9e1bcfed71e6951b953611d11dda32ea09d78","0x78295e5a22b84e982cf601eb639597b8b0515a88cb5ac7fa8a4aabe3c87349d","0x2fa5e5f18f6027a6501bec864564472a616b2e274a41211a444cbe3a99f3cc61","0xe884376d0d8fd21ecb780389e941f66e45e7acce3e228ab3e2156a614fcd747","0x1b7201da72494f1e28717ad1a52eb469f95892f957713533de6175e5da190af2","0x1f8d8822725e36385200c0b201249819a6e6e1e4650808b5bebc6bface7d7636","0x2c5d82f66c914bafb9701589ba8cfcfb6162b0a12acf88a8d0879a0471b5f85a","0x14c54148a0940bb820957f5adf3fa1134ef5c4aaa113f4646458f270e0bfbfd0","0x190d33b12f986f961e10c0ee44d8b9af11be25588cad89d416118e4bf4ebe80c","0x22f98aa9ce704152ac17354914ad73ed1167ae6596af510aa5b3649325e06c92","0x2a7c7c9b6ce5880b9f6f228d72bf6a575a526f29c66ecceef8b753d38bba7323","0x2e8186e558698ec1c67af9c14d463ffc470043c9c2988b954d75dd643f36b992","0xf57c5571e9a4eab49e2c8cf050dae948aef6ead647392273546249d1c1ff10f","0x1830ee67b5fb554ad5f63d4388800e1cfe78e310697d46e43c9ce36134f72cca","0x2134e76ac5d21aab186c2be1dd8f84ee880a1e46eaf712f9d371b6df22191f3e","0x19df90ec844ebc4ffeebd866f33859b0c051d8c958ee3aa88f8f8df3db91a5b1"]]}`
	response, err := http.Post("http://localhost:8080/prove", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}
}

// func TestInsertionHappyPath(t *testing.T) {
// 	if mode != server.InsertionMode {
// 		return
// 	}
// 	body := `{
// 		"inputHash":"0x5057a31740d54d42ac70c05e0768fb770c682cb2c559bdd03fe4099f7e584e4f",
// 		"startIndex":0,
// 		"preRoot":"0x18f43331537ee2af2e3d758d50f72106467c6eea50371dd528d57eb2b856d238",
// 		"postRoot":"0x2267bee7aae8ed55eb9aecff101145335ed1dd0a5a276a2b7eb3ae7d20e232d8",
// 		"identityCommitments":["0x1","0x2"],
// 		"merkleProofs": [
// 			["0x0","0x2098f5fb9e239eab3ceac3f27b81e481dc3124d55ffed523a839ee8446b64864","0x1069673dcdb12263df301a6ff584a7ec261a44cb9dc68df067a4774460b1f1e1"],
// 			["0x1","0x2098f5fb9e239eab3ceac3f27b81e481dc3124d55ffed523a839ee8446b64864","0x1069673dcdb12263df301a6ff584a7ec261a44cb9dc68df067a4774460b1f1e1"]
// 		]}`
// 	response, err := http.Post("http://localhost:8080/prove", "application/json", strings.NewReader(body))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if response.StatusCode != http.StatusOK {
// 		t.Fatalf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
// 	}
// }

// func TestDeletionHappyPath(t *testing.T) {
// 	if mode != server.DeletionMode {
// 		return
// 	}
// 	body := `{
// 		"inputHash":"0xdcd389a94b549222fadc9e335c358a3fe4d534155182f46927f82ea8491c7480",
// 		"deletionIndices":[0,2],
// 		"preRoot":"0xd11eefe87b985333c0d327b0cdd39a9641b5ac32c35c2bda84301ef3231a8ac",
// 		"postRoot":"0x1912415186579e1d9ff6282b76d081f0acd527d8549ea803385b1382d9498f35",
// 		"identityCommitments":["0x1","0x3"],
// 		"merkleProofs":[
// 			["0x2","0x20a3af0435914ccd84b806164531b0cd36e37d4efb93efab76913a93e1f30996","0x1069673dcdb12263df301a6ff584a7ec261a44cb9dc68df067a4774460b1f1e1"],
// 			["0x4","0x65e2c6cc08a36c4a943286bc91c216054a1981eb4f7570f67394ef8937a21b8","0x1069673dcdb12263df301a6ff584a7ec261a44cb9dc68df067a4774460b1f1e1"]
// 		]}`
// 	response, err := http.Post("http://localhost:8080/prove", "application/json", strings.NewReader(body))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if response.StatusCode != http.StatusOK {
// 		t.Fatalf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
// 	}
// }

// func TestInsertionWrongInput(t *testing.T) {
// 	if mode != server.InsertionMode {
// 		return
// 	}
// 	body := `{
// 		"inputHash":"0x5057a31740d54d42ac70c05e0768fb770c682cb2c559bdd03fe4099f7e584e4f",
// 		"startIndex":0,
// 		"preRoot":"0x18f43331537ee2af2e3d758d50f72106467c6eea50371dd528d57eb2b856d238",
// 		"postRoot":"0x2267bee7aae8ed55eb9aecff101145335ed1dd0a5a276a2b7eb3ae7d20e232d8",
// 		"identityCommitments":["0x1","0x2"],
// 		"merkleProofs": [
// 			["0x0","0x0","0x0"],
// 			["0x1","0x0","0x0"]
// 		]}`
// 	response, err := http.Post("http://localhost:8080/prove", "application/json", strings.NewReader(body))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if response.StatusCode != http.StatusBadRequest {
// 		t.Fatalf("Expected status code %d, got %d", http.StatusBadRequest, response.StatusCode)
// 	}
// 	responseBody, err := io.ReadAll(response.Body)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if !strings.Contains(string(responseBody), "proving_error") {
// 		t.Fatalf("Expected error message to be tagged with 'proving_error', got %s", string(responseBody))
// 	}

// }

// func TestDeletionWrongInput(t *testing.T) {
// 	if mode != server.DeletionMode {
// 		return
// 	}
// 	body := `{
// 		"inputHash":"0xdcd389a94b549222fadc9e335c358a3fe4d534155182f46927f82ea8491c7480",
// 		"deletionIndices":[0,2],
// 		"preRoot":"0xd11eefe87b985333c0d327b0cdd39a9641b5ac32c35c2bda84301ef3231a8ac",
// 		"postRoot":"0x1912415186579e1d9ff6282b76d081f0acd527d8549ea803385b1382d9498f35",
// 		"identityCommitments":["0x1","0x3"],
// 		"merkleProofs":[
// 			["0x2","0xD","0xD"],
// 			["0x4","0xD","0xD"]
// 		]}`
// 	response, err := http.Post("http://localhost:8080/prove", "application/json", strings.NewReader(body))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if response.StatusCode != http.StatusBadRequest {
// 		t.Fatalf("Expected status code %d, got %d", http.StatusBadRequest, response.StatusCode)
// 	}
// 	responseBody, err := io.ReadAll(response.Body)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if !strings.Contains(string(responseBody), "proving_error") {
// 		t.Fatalf("Expected error message to be tagged with 'proving_error', got %s", string(responseBody))
// 	}
// }

// func TestDeletionBatchPadding(t *testing.T) {
// 	if mode != server.DeletionMode {
// 		return
// 	}
// 	body := `{
// 		"inputHash":"0x509d6e4ca8a621713cc5feb95de95cb4eed3c1127176d93da653fd3cc55db537",
// 		"deletionIndices":[0,8],
// 		"preRoot":"0xd11eefe87b985333c0d327b0cdd39a9641b5ac32c35c2bda84301ef3231a8ac",
// 		"postRoot":"0x22c58cf24838c2eb1701f2aa6e6a867e10237590dbdb423e4d3e053b121c44cb",
// 		"identityCommitments":["0x1","0x0"],
// 		"merkleProofs":[
// 			["0x2","0x20a3af0435914ccd84b806164531b0cd36e37d4efb93efab76913a93e1f30996","0x1069673dcdb12263df301a6ff584a7ec261a44cb9dc68df067a4774460b1f1e1"],
// 			["0x0","0x0","0x0"]
// 		]}`
// 	response, err := http.Post("http://localhost:8080/prove", "application/json", strings.NewReader(body))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if response.StatusCode != http.StatusOK {
// 		t.Fatalf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
// 	}
// }
