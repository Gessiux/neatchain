// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package params

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the INT Chain main network.
var MainnetBootnodes = []string{

	"enode://e1fadafa76bad701b108e03e4781cc6cf65df32e9773a7b7175d3b542986baade931c186374b07eaedf31a4d65e2aecabc28fcbb8a6f586110c436198c59d192@135.181.195.79:8550",
}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// INT Chain test network.
var TestnetBootnodes = []string{
	"enode://20c9a43afb9b19ada68883ef404810d93963be0e09f90444cce52e9a68fbdb9b7fec0187dd1e2c31737c0b452884a43e0fd654bde18310588b368d080836f8f9@101.32.74.50:8550", // Titans
	//"enode://6eb4941356e557692057bd58c6324ccb6fff2864389e73e67d661a83114d4e9218c1a00b018d38ef4076f2581ecc48aa5e4ec37da241c42ef4817532dcc4ad0b@129.226.59.148:8551", // Oceanus
	//"enode://2c881d1e8eb0516e4b203aa991df9f2ea5ec2a962c54b1a383451f7e6b83e729c75b819ecbd3946e847fee18c2915d0a1cc0425871c1535b7449a4809973c5ee@129.226.128.55:8551", // Iapetus
	//"enode://d8f5598499a106b48d50080637b9ee1b01d2cea2c545348508a6b3310a9d77c17fca253cb1bad3c5bd4e6cc7aae72cc5c01a099458a4f428ecc50469dc760168@129.226.63.13:8551",  // Mnemosyne
	//"enode://bc51a4ca30d02ec7a7926dade9d1a65dfb3a5e73c70e600ada29a51e8cdfdee79c81b1cdc12b2a5d001a93329bdd7a434592d6dd9b152739c001839ce7095c6d@170.106.160.155:8551",
	//"enode://7ebef9b823797285be19487a8a7d237735d0c6a7546761d720ab67c13799d47bb69a5742a39f4cb5804bb246bba2757c3afdce580d5932256ba66e367e12bb54@170.106.9.165:8551",
}
