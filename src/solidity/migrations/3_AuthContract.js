const AuthContract = artifacts.require("AuthContract");

module.exports = function (deployer) {
  deployer.deploy(AuthContract, '0x0f73f1c6c755eb82ef1c494d1afd455c10f00cf3', '0x0f73f1c6c755eb82ef1c494d1afd455410f00cf3');
};