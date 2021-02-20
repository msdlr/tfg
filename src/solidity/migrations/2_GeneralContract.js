const GeneralContract = artifacts.require("GeneralContract");

module.exports = function (deployer) {
  deployer.deploy(GeneralContract);
};
