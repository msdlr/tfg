const GeneralContract = artifacts.require("GeneralContract");

module.exports = function (deployer) {
  deployer.deploy(GeneralContract, '0x0f73f1c6c755eb82ef1c494d1afd455c10f00cf3', 'dajshdgasd');
};
