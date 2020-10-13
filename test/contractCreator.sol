pragma solidity >=0.4.22 <0.7.0;
import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/GeneralContract.sol";

contract contractCreator {
    
    GeneralContract testContract;
    address thisContract = address(this);


    // Create a contract with THIS contract as admin so that the other contract
    // is just a regular user
    function createContract() public returns (GeneralContract){
        testContract = new GeneralContract(thisContract, "abc");
        return testContract;
    }

    // Add the caller (another contract) to the user list
    function addUser() public {
        testContract.addUser(msg.sender,"user1");
    }
}