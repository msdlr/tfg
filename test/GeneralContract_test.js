const GeneralContract = artifacts.require("GeneralContract");

contract('GeneralContract', function(accounts){

    it("should create the contract",function(){
        return GeneralContract.deployed().then(function(instance){
            return instance.getOwner.call();
        }).then(function(getOwner){
            assert.equal(getOwner.valueOf(), '0x0f73f1c6c755eb82ef1c494d1afd455c10f00cf3', "ADSDASD");
        })
    })

});