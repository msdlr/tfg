pragma solidity >=0.4.22 <0.7.0;
import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/GeneralContract.sol";

// File name has to end with '_test.sol', this file can contain more than one testSuite contracts
contract GeneralContract_Admintest {
    
    GeneralContract testContract;
    address thisContract = address(this);
    
    address[] testAddrs = 
    [0xAb8483F64d9C6d1EcF9b849Ae677dD3315835cb2,
    0x4B20993Bc481177ec7E8f571ceCaE8A9e22C02db,
    0x4B0897b0513fdC7C541B6d9D7E929C4e5364D2dB,
    0x583031D1113aD414F02576BD6afaBfb302140225,
    0xdD870fA1b7C4700F2BD7f44238821C26f7392148];
    

    /// 'beforeAll' runs before all other tests
    /// More special functions are: 'beforeEach', 'beforeAll', 'afterEach' & 'afterAll'

    // msg.sender in GeneralContract -> address (this) in this contract
    function test_createContract() public {
        // Instantiate the contract to test
        testContract = new GeneralContract(thisContract, "M4573R");
        // Check that the constructor executed correctly
        //Assert.equal(thisContract, testContract.getOwner(), "owner address should be this caller");
        
        // Check if the fields initialized correctly
        Assert.isTrue(testContract.getUserAdmin(thisContract),"caller was not made admin");
        Assert.isTrue(testContract.getUserRegistered(thisContract),"caller is not on the user list");
        Assert.isFalse(testContract.getUserLoggedIn(thisContract),"caller is logged in");
        Assert.equal("M4573R", testContract.getUserId(thisContract), "contract id was not set");
    }

    function test_addUser() public {
        
        // Try to add a new user
        //Assert.isFalse(testContract.getUserRegistered(testAddrs[0]),"User entry (isRegistered) should be FALSE");
        testContract.addUser(testAddrs[0], "user0");
        
        // Check it was done correctly
        Assert.isTrue(testContract.getUserRegistered(testAddrs[0]),"User entry (isRegistered) should be TRUE");
    }

    function test_promoteUser() public {
        //Assert.isFalse(testContract.getUserAdmin(testAddrs[0]),"User entry (adminStatus) should be FALSE");
        testContract.promoteUser(testAddrs[0]);
        Assert.isTrue(testContract.getUserAdmin(testAddrs[0]),"User entry (adminStatus) should be FALSE");
    }

    function test_demoteUser() public {
        //Assert.isTrue(testContract.getUserAdmin(testAddrs[0]),"User entry (adminStatus) should be FALSE");
        testContract.demoteAdmin(testAddrs[0]);
        Assert.isFalse(testContract.getUserAdmin(testAddrs[0]),"User entry (adminStatus) should be FALSE");
    }
    
    function test_removeUser() public {
        //Assert.isTrue(testContract.getUserRegistered(testAddrs[0]),"User to test is not registered");
        //Assert.equal("user0",testContract.getUserId(testAddrs[0]),"User to test doesn't match id provided");
        testContract.rmUser(testAddrs[0],"user0");

        Assert.isFalse(testContract.getUserRegistered(testAddrs[0]),"user was not removed");
        Assert.equal("",testContract.getUserId(testAddrs[0]),"User struct still contains id");
        Assert.equal(testContract.getUserAddress("user0"),address(0),"Id index still points to user address");
        Assert.isFalse(testContract.getUserAdmin(testAddrs[0]), "adminStatus field not cleared after removal");
        Assert.isFalse(testContract.getUserLoggedIn(testAddrs[0]), "isLoggedIn field not cleared after removal");
        Assert.equal(uint(0),uint(testContract.getUserAttempts(testAddrs[0])), "attempts field not cleared after removal");
    }
    

    
}